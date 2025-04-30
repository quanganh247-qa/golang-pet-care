# Single WebSocket Connection Guide

This document explains how to use the new WebSocket system with a single connection.

## Connection

Connect to the WebSocket endpoint without any token:

```javascript
const socket = new WebSocket('ws://your-api-domain/ws/chat');

socket.onopen = () => {
  console.log('WebSocket connection established');
};

socket.onmessage = (event) => {
  const message = JSON.parse(event.data);
  handleWebSocketMessage(message);
};

socket.onerror = (error) => {
  console.error('WebSocket error:', error);
};

socket.onclose = () => {
  console.log('WebSocket connection closed');
};
```

## Authentication

After connecting, you need to authenticate. You can do this in two ways:

### 1. Using Username and Password

```javascript
function authenticate(username, password) {
  const authMessage = {
    type: 'authenticate',
    data: {
      username: username,
      password: password
    }
  };
  socket.send(JSON.stringify(authMessage));
}
```

### 2. Using a Token

```javascript
function authenticateWithToken(token) {
  const authMessage = {
    type: 'authenticate',
    data: {
      token: token
    }
  };
  socket.send(JSON.stringify(authMessage));
}
```

## Handling Messages

Handle different types of WebSocket messages:

```javascript
function handleWebSocketMessage(message) {
  switch (message.type) {
    case 'connected':
      console.log('Connected to WebSocket server');
      break;
      
    case 'auth_required':
      console.log('Authentication required');
      // Trigger authentication process here
      break;
      
    case 'auth_success':
      console.log('Authentication successful');
      console.log('User ID:', message.data.user_id);
      console.log('Username:', message.data.username);
      // Store authentication state
      // Start heartbeat
      startHeartbeat();
      break;
      
    case 'auth_error':
      console.error('Authentication error:', message.message);
      break;
      
    case 'chat_message':
      // Handle chat message
      console.log('Received chat message:', message.data);
      break;
      
    case 'typing':
      // Handle typing indication
      console.log('User is typing:', message.data);
      break;
      
    case 'error':
      console.error('Error:', message.message);
      break;
  }
}
```

## Sending Messages

Send messages after authentication:

```javascript
function sendChatMessage(conversationId, content) {
  const chatMessage = {
    type: 'chat_message',
    conversationId: conversationId,
    data: {
      content: content,
      type: 'text'
    }
  };
  socket.send(JSON.stringify(chatMessage));
}

function sendTypingStatus(conversationId) {
  const typingMessage = {
    type: 'typing',
    conversationId: conversationId
  };
  socket.send(JSON.stringify(typingMessage));
}

function sendReadReceipt(messageId) {
  const readMessage = {
    type: 'read_receipt',
    messageId: messageId
  };
  socket.send(JSON.stringify(readMessage));
}
```

## Heartbeat

Keep the connection alive:

```javascript
function startHeartbeat() {
  setInterval(() => {
    if (socket.readyState === WebSocket.OPEN) {
      const heartbeat = {
        type: 'heartbeat'
      };
      socket.send(JSON.stringify(heartbeat));
    }
  }, 30000); // 30 seconds
}
```

## Reconnection Strategy

Implement reconnection logic:

```javascript
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;
const reconnectDelay = 3000; // 3 seconds

function connectWebSocket() {
  socket = new WebSocket('ws://your-api-domain/ws/chat');
  
  socket.onopen = () => {
    console.log('WebSocket connection established');
    reconnectAttempts = 0;
    // If you were authenticated before, re-authenticate
    if (authToken) {
      authenticateWithToken(authToken);
    }
  };
  
  socket.onclose = (event) => {
    if (!event.wasClean && reconnectAttempts < maxReconnectAttempts) {
      reconnectAttempts++;
      console.log(`Attempting to reconnect (${reconnectAttempts}/${maxReconnectAttempts})...`);
      setTimeout(connectWebSocket, reconnectDelay);
    } else {
      console.log('WebSocket connection closed');
    }
  };
  
  // Other event handlers...
}

// Initial connection
connectWebSocket();
``` 