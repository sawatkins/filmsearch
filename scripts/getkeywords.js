require('dotenv').config(); 
const fs = require('fs');
const axios = require('axios');

const API_KEY = process.env.TMDB_API_KEY;

let requestCount = 0;

const fetchKeywords = async (movieId) => {
  const url = `https://api.themoviedb.org/3/movie/${movieId}/keywords?api_key=${API_KEY}`;
  const response = await axios.get(url);
  return response.data.keywords.map(keyword => keyword.name);
};

const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms));

const main = async () => {
  try {
    const fileData = fs.readFileSync('cleantopmovies.json', 'utf8');
    const parsedData = JSON.parse(fileData);

    for (let i = 0; i < parsedData.results.length; i++) {
      const movie = parsedData.results[i];

      if (requestCount >= 40) {
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

    fs.writeFileSync('topmovieswkeywords.json', JSON.stringify(parsedData, null, 2));
    console.log('Successfully written to topmovieswkeywords.json');
  } catch (error) {
    console.error('Error:', error);
  }
};

main();
