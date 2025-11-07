const Actions = {
 SendMessage:1,
	SendAllCli : 2,
	GetAllCli : 3,
	CreateCli   : 4,

	GetKey   : 8,
	RegistryKey : 9,
	UpdateKey   : 10,
	NeedKey     : 12,
 
	Disconnected :11,
};
async function importEcdsaP256PublicKeyJwk(jwk) {
  const publicKey = await window.crypto.subtle.importKey(
    'jwk',
    jwk,
    {
      name: 'ECDSA',
      namedCurve: 'P-256',
    },
    true, // Public keys are generally extractable
    ['verify']
  );

  return publicKey;
}
class Clients {
  constructor(){
    this.clients = new Map();

  }
  setUser(login, shared,identityKey){
    this.clients.set(login, { shared,identityKey, limit: 50 });
  }
  deleteUser(login){
    this.clients.delete(login)
  }
  isUser(login){
    return this.clients.has(login)
  }
  getIdentity(login){
    return this.clients.get(login).identityKey
  }
   getUser(login) {
    const entry = this.clients.get(login);
    if (!entry) return null;

    entry.limit -= 1;

    // если лимит кончился — удаляем и возвращаем null (чтобы запустить новый ратчет)
    if (entry.limit <= 0) {
      this.clients.delete(login);
      return entry.shared;
    }

    return entry.shared;
  }
}
const clientes = new Clients
//crypted logic
function getRandomIv() {
  return crypto.getRandomValues(new Uint8Array(12));
}
async function encryptMessage(sharedSecret, message) {
  const enc = new TextEncoder();
  const iv = getRandomIv();
  const ciphertext = await crypto.subtle.encrypt(
    {
      name: "AES-GCM",
      iv: iv,
    },
    sharedSecret,
    enc.encode(message)
  );
  return { ciphertext, iv };
}
//finish

//key storage
class KeyStore {
  constructor() {
    this.identityKey = null;  // CryptoKey приватный identity
    this.preKeys = [];        // массив {publicKey, privateKey}
    this.currIndex = 0;       // для round-robin
  }

  async addIdentityKey(keyPair) {
    if (this.identityKey!=null){
      return
    }
    this.identityKey = keyPair;
  }
  async getIdentityKey(){
    return this.identityKey
  }
  async addPreKeys(keys) {
    this.preKeys = keys;
    this.currIndex = 0;
  }

  // Берем следующий ключ для шифрования (round-robin)
  getNextPreKey() {
    if (!this.preKeys.length) throw new Error("No prekeys available");
    const key = this.preKeys[this.currIndex];
    this.currIndex = (this.currIndex + 1) % this.preKeys.length;
    return key;
  }

  // Поиск приватного ключа по публичному (для расшифровки)
  findPrivateKeyByPublic(pubKey) {
    const k = this.preKeys.find(k => k.publicKey === pubKey);
    if (!k) throw new Error("Key not found");
    return k.privateKey;
  }
}
//finish


///keys logic
// cryptoKey — это объект типа CryptoKey (public)
 async function exportPublicKeyII(publicKey) {
            // Export as SPKI (PEM format)
            const spkiKey = await crypto.subtle.exportKey("spki", publicKey);
            const spkiPem = `-----BEGIN PUBLIC KEY-----\n${btoa(String.fromCharCode(...new Uint8Array(spkiKey)))}\n-----END PUBLIC KEY-----`;
           

            // Export as JWK
            const jwkPublicKey = await crypto.subtle.exportKey("jwk", publicKey);
            return jwkPublicKey
        }

async function generateKeyPair() {
    return crypto.subtle.generateKey(
        {
            name: "ECDH",
            namedCurve: "P-256",
        },
        true,
        ["deriveKey", "deriveBits"]
    );
}
async function generateIdentityKey() {
  const keyPair = await crypto.subtle.generateKey(
    { name: "ECDSA", namedCurve: "P-256" },
    true,
    ["sign", "verify"]
  );
  return keyPair; // {publicKey, privateKey}
}

//из b64 получателя делаем объект крипто кей
async function importPublicKey(b64Key) {
  // Конвертируем base64 в ArrayBuffer
  const binary = atob(b64Key);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i);
  }

  // Импортируем ключ в формате raw для ECDH
  return crypto.subtle.importKey(
    "raw",               // формат
    bytes.buffer,        // ArrayBuffer
    { name: "ECDH", namedCurve: "P-256" },
    true,                // extractable
    []                   // usages для публичного ключа не нужны
  );
}

async function exportPublicKey(keyPair) {
    const raw = await crypto.subtle.exportKey("raw", keyPair.publicKey);
    return btoa(String.fromCharCode(...new Uint8Array(raw)));
}

async function generatePrekeys(count) {
    const prekeys = [];
    for (let i = 0; i < count; i++) {
        const kp = await generateKeyPair();
        const pub = await exportPublicKey(kp);
        prekeys.push({ publicKey: pub, privateKey: kp.privateKey });
    }
    return prekeys;
}
//keys logic finish

//base64 logic
function toBase64(txt)
{
    // TextEncoder: Always UTF8
    const uint8Array = new TextEncoder().encode(txt);
    let binary = '';

    for (let i = 0; i < uint8Array.length; ++i)
        binary += String.fromCharCode(uint8Array[i]);

    return btoa(binary);
}

//finish

//refresh users from server logic
function getUser(){
   const msg = {
    Action: 3,
    Message: messageInput.value,
   // To: selectedType === 'private' ? prompt("Enter recipient username") : '',
    From: username, 
  };
  const jMsg = JSON.stringify(msg);
 const b64Msg = toBase64(jMsg);
  socket.send(b64Msg);
}
document.getElementById('refreshButton').addEventListener('click', getUser);
function userPicker(json){
 if (json.Action === Actions.GetAllCli) {
  let users;

  // Проверяем, если Message — это строка, пытаемся распарсить её
  if (typeof json.Message === 'string') {
    try {
      users = JSON.parse(json.Message);
    } catch (e) {
      console.error('Ошибка парсинга Message:', e);
      users = [];
    }
  } else if (Array.isArray(json.Message)) {
    users = json.Message;
  } else {
    users = [];
  }

  // Теперь users — это реальный массив
  if (Array.isArray(users)) {
    

    const userPicker = document.getElementById('userSelect');
    userPicker.innerHTML = ''; // очищаем старые опции

    users.forEach(user => {
      const option = document.createElement('option');
      option.value = user;
      option.textContent = user;
      userPicker.appendChild(option);
    });
  }
  
}

}
//finish

function arrayBufferToBase64(buffer) {
  const bytes = new Uint8Array(buffer);
  let binary = '';
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary);
}

// Establish the WebSocket connection
// Use 'wss://' for secure connections (recommended for production)
// Use 'ws://' for unsecure connections (e.g., local development)
const socket = new WebSocket('ws://localhost:8080/ws'); 
let username
// Event listener for when the connection is successfully established
socket.onopen = (event) => {
  console.log('WebSocket connection established:', event);
  socket.send(username=prompt("entr you username"))
  document.getElementById('mylogin').textContent='my login: ' +username
  getUser();
};
const keyStore = new KeyStore();
// message from server to cli
socket.onmessage = (event) => {
  
  // Process the received data
  // Example: Display the message on the webpage
  const messagesDiv = document.getElementById('messages');
  if (messagesDiv) {
   unb64= atob(event.data)
   const decoded = new TextDecoder().decode(
  Uint8Array.from(unb64, c => c.charCodeAt(0))
);
     const json = JSON.parse(decoded);
  
   const messageText = json.Message;
   const messageFrom = json.From;
   const dateFrom = dateFormatter(json.Datatime);
   const fullMessage = dateFrom + ' ' + messageFrom + ': ' + messageText
    switch (json.Action){
      case Actions.SendMessage:
          (async()=>{
            
             const shared = clientes.getUser(json.From)
             const key = clientes.getIdentity(json.From)
             const iv = Uint8Array.from(atob(json.Keys.Iv), c => c.charCodeAt(0));
            const ciphertext = Uint8Array.from(atob(json.Message), c => c.charCodeAt(0));
           const signature = Uint8Array.from(atob(json.Keys.Signature), c => c.charCodeAt(0));
            
             const isValid = await crypto.subtle.verify(
                    {
                      name: "ECDSA",
                      hash: { name: "SHA-256" },
                    },
                    key,
                    signature,
                    ciphertext
                  );
              if (isValid){
                  const decrypted = await crypto.subtle.decrypt(
                  {
                    name: "AES-GCM",
                    iv: iv
                  },
                  shared,
                  ciphertext
                );
               const decodedMessage = new TextDecoder().decode(decrypted);
                 
                messagesDiv.innerHTML += `<p style="color: purple;">${dateFrom+'⬅️'+'['+json.From+']' + ': ' + decodedMessage}</p>`;

              }
             
              })()
            break
      case Actions.SendAllCli: 
         messagesDiv.innerHTML += `<p>${fullMessage}</p>`;
         break
      case Actions.GetAllCli:
        userPicker(json)
        break
      case Actions.Disconnected:
        clientes.deleteUser(json.Message)
        messagesDiv.innerHTML += `<p style="color: red;">${'['+json.Message+']'+ ' disconnected'}</p>`;
        break
        //просто отправляем ключи серверу
      case Actions.NeedKey:
          (async () => {
            const identityKey = await generateIdentityKey();
            const preKeys = await generatePrekeys(10);
            const publicjwt = await exportPublicKeyII(identityKey.publicKey);
            
            const publicB64 = btoa(JSON.stringify(publicjwt));
            
            keyStore.addIdentityKey(identityKey);
            keyStore.addPreKeys(preKeys);
            // Создаем объект сообщения
            const msg = {
              Action: Actions.RegistryKey,
              From: username,
              Keys: {
                Identity: publicB64,
                Keys: preKeys.map(k => k.publicKey),
              }
            };
            
            // Сериализуем в JSON
           const jsonStr = JSON.stringify(msg);
            const b64msg = toBase64(jsonStr);
            socket.send(b64msg);
            console.log('keys sent server')
          })();
          break;
          case Actions.GetKey:
          
            (async()=>{
              const peerKeyB64 = json.Keys.Key; // ключ клиента Б
              const peerKey = await importPublicKey(peerKeyB64);
              const myCryptoKey = keyStore.getNextPreKey(); // пример метода
              const sharedSecret = await crypto.subtle.deriveKey(
                              {
                                name: "ECDH",
                                public: peerKey,             // публичный ключ получателя
                              },
                              myCryptoKey.privateKey,                 // ваш приватный ключ
                              {
                                name: "AES-GCM",            // для шифрования сообщений
                                length: 256,
                              },
                              false,                        // экспортировать нельзя (или true если нужно)
                              ["encrypt", "decrypt"]
                            );
                const senderIdentityKeyb64 = json.Keys.Identity;
                const jwkString= atob(senderIdentityKeyb64) 
                const jwk = JSON.parse(jwkString);
                const key = await importEcdsaP256PublicKeyJwk(jwk)
                clientes.setUser(json.To,sharedSecret,key)
               console.log('key exchange done')
                // messagesDiv.innerHTML += `<p style="color: green;">${'обмен ключами шифрования выполнен c '+ json.To}</p>`;

            })(); 
      }
  }
};

// Event listener for errors during the connection
socket.onerror = (error) => {
  console.error('WebSocket error:', error);
};

// Event listener for when the connection is closed
socket.onclose = (event) => {
  if (event.wasClean) {
    console.log(`Connection closed cleanly, code=${event.code}, reason=${event.reason}`);
  } else {
    console.error('Connection died unexpectedly');
  }
};


function dateFormatter(datatime){
  const date = new Date(datatime);
  return '['+date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })+']';
}


// Send message from cli to sever
function sendMessage() {
  const messagesDiv = document.getElementById('messages');
  const messageInput = document.getElementById('messageInput');
  const msgTypeSelect = document.getElementById('msgType');

  if (!messageInput || socket.readyState !== WebSocket.OPEN) {
    console.warn('WebSocket is not open. Cannot send message.');
    return;
  }
 const  userTo = document.getElementById('userSelect').value;
  // Определяем Action по выбору
  
  const selectedType = msgTypeSelect.value;
  let msg ={}
  switch (selectedType) {
    case 'public':
             msg = {
              Action: Actions.SendAllCli,
              Message: messageInput.value,
              Datatime: Date.now(),
              To: userTo,
              From: username, // можно заполнять логином
            }; 
      messageInput.value = '';
      break;
    case 'private':
      (async()=>{
        let isUser = clientes.isUser(userTo)
        if (!isUser){
          msg = {
            Action: Actions.GetKey,
            Message: "",
            Datatime: Date.now(),
            To: userTo,
            From: username
          };
           const jMsg = JSON.stringify(msg);
           const b64Msg = toBase64(jMsg);
           socket.send(b64Msg);
        }
        const start = Date.now();
        while (!isUser && Date.now() - start < 10000) {
        await new Promise(r => setTimeout(r, 200));
        isUser = clientes.isUser(userTo);
      }
        if (!isUser) {
        messagesDiv.innerHTML += `<p style="color: red;">⚠️ Не удалось получить общий ключ с ${userTo}</p>`;
        return;
      }
      const secr = clientes.getUser(userTo)
      const { ciphertext, iv } = await encryptMessage(secr, messageInput.value);
      const signature = await crypto.subtle.sign(
                  {
                    name: "ECDSA",
                    hash: { name: "SHA-256" },
                  },
                  keyStore.identityKey.privateKey,  // приватный identity key отправителя
                  ciphertext
                );
            msg = {
            Action: Actions.SendMessage,
            From: username,
            To: userTo,
            Datatime: Date.now(),
            Message: btoa(String.fromCharCode(...new Uint8Array(ciphertext))),
            Keys: {
              Iv: btoa(String.fromCharCode(...iv)),
              Signature: arrayBufferToBase64(signature),
            }
          };
          const jMsg = JSON.stringify(msg);
           const b64Msg = toBase64(jMsg);

             socket.send(b64Msg);
          messagesDiv.innerHTML += `<p style="color: purple;">${dateFormatter(msg.Datatime)+'➡️' + '['+userTo+']' + ': ' + messageInput.value}</p>`;

          messageInput.value = ''
      })()
        break;
    default:
       messageInput.value = '';
  }

  const jMsg = JSON.stringify(msg);
  const b64Msg = toBase64(jMsg);

  socket.send(b64Msg);

 


}

// Привязка к кнопке
document.getElementById('sendButton').addEventListener('click', sendMessage);
