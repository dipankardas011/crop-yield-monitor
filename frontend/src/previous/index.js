// index.js

const express = require('express');
const app = express();
const port = 3000; // Choose your desired port

// Serve static files (styles and scripts)
app.use(express.static('public'));

// Define routes
app.get('/', (req, res) => {
    res.sendFile(__dirname + '/public/index.html');
});

// Start the server
app.listen(port, () => {
    console.log(`Server is running at http://localhost:${port}`);
});
