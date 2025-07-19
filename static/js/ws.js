const messageContainer = document.getElementById('messages');
const websocketURL = `${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${location.host}/ws/message-stream`;
const messageURL = `/post-message`;

// WebSocket connection handling
function connect() {
	const ws = new WebSocket(websocketURL);

	ws.onopen = () => {
		console.log("Connected to WebSocket server");
		// Could update UI to show connected status
	};

	ws.onmessage = (event) => {
		try {
			const message = JSON.parse(event.data);
			if (message.type === "message" || message.type === "connection") {
				const decodedPayload = atob(message.payload);
				const newMessage = document.createElement('p');
				newMessage.textContent = decodedPayload;
				messageContainer.appendChild(newMessage);
			} else {
				console.log("Received a different type of message:", message);
			}
		} catch (error) {
			console.error("Error processing message:", error);
			console.log("Raw message:", event.data);
		}
	};

	ws.onclose = () => {
		console.log("Disconnected from WebSocket server");
		// Attempt to reconnect after 1 second
		setTimeout(() => {
			console.log("Attempting to reconnect...");
			websocket = connect();
		}, 1000);
	};

	ws.onerror = (error) => {
		console.error("WebSocket error:", error);
	};

	return ws;
}

// Initialize WebSocket connection
let websocket = connect();

// HTTP message sending
function sendMessage() {
	const messageInput = document.getElementById('messageInput');
	const message = messageInput.value;

	if (message.trim() === "") {
		alert("Please enter a message.");
		return;
	}

	fetch(messageURL, {
		method: 'POST',
		headers: {
			'message': message
		},
	})
		.then(response => {
			if (!response.ok) {
				return response.text().then(err => { throw new Error(err || 'HTTP error') });
			}
			return response.text();
		})
		.then(() => {
			console.log('Message sent successfully');
			messageInput.value = '';
		})
		.catch(error => {
			console.error('Error sending message:', error);
		});
}
