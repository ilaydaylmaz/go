const express = require('express');
const sqlite3 = require('sqlite3').verbose();
const app = express();
const port = 3001;

let db = new sqlite3.Database('./user.db', (err) => {
 if (err) {
    return console.error(err.message);
 }
 console.log('Connected to the SQLite database.');
});

app.use(express.json());

app.post('/add-user', (req, res) => {
 const { name, email } = req.body;

 db.run('INSERT INTO users (name, email) VALUES (?, ?)', [name, email], function (err) {
    if (err) {
      return console.log(err.message);
    }
    console.log(`User added with ID: ${this.lastID}`);
    res.status(200).json({ message: 'User added successfully.' });
 });
});

app.listen(port, () => {
 console.log(`API listening at http://localhost:${port}`);
});