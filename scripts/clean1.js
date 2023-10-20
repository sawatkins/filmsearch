const fs = require('fs');

// Read the input file
const data = JSON.parse(fs.readFileSync('topmovies.json'));

// Extract the id and title for each item
const results = [];
for (const item of data.results) {
    results.push({id: item.id, title: item.title});
}

// Write the output file
fs.writeFileSync('cleantopmovies.json', JSON.stringify({results: results}));
