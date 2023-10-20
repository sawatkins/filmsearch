const fs = require('fs');
const axios = require('axios');

const API_KEY = 'a19c17e1944a778874b2e9598b42b5ff'; 

let requestCount = 0;

// Function to fetch keywords for a movie by its ID
const fetchKeywords = async (movieId) => {
  const url = `https://api.themoviedb.org/3/movie/${movieId}/keywords?api_key=${API_KEY}`;
  const response = await axios.get(url);
  return response.data.keywords.map(keyword => keyword.name);
};

// Sleep function
const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms));

// Main function to read JSON file, fetch keywords, and write to a new JSON file
const main = async () => {
  try {
    // Read the JSON file
    const fileData = fs.readFileSync('cleantopmovies.json', 'utf8');
    const parsedData = JSON.parse(fileData);

    for (let i = 0; i < parsedData.results.length; i++) {
      const movie = parsedData.results[i];

      if (requestCount >= 40) {
        // Sleep for 2 seconds if 40 requests have been made
        console.log('Rate limit reached, sleeping for 2 seconds');
        await sleep(2000);
        requestCount = 0;
      }

      console.log("fetching movie:" + movie.title)
      const keywords = await fetchKeywords(movie.id);
      console.log(keywords)
      movie['keywords'] = keywords;

      requestCount++;
    }

    // Write to a new JSON file
    fs.writeFileSync('topmovieswkeywords.json', JSON.stringify(parsedData, null, 2));
    console.log('Successfully written to topmovieswkeywords.json');
  } catch (error) {
    console.error('Error:', error);
  }
};

// Run the main function
main();
