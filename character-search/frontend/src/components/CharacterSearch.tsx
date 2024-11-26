// 

import React, { useState } from "react";
import axios from "axios";

interface Character {
  _id: number;
  name: string;
  status: string;
  species: string;
  gender: string;
  image: string;
}

const CharacterSearch = () => {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<Character[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const searchCharacters = async () => {
    setLoading(true);
    setError("");
    try {
      const response = await axios.get(`http://localhost:8080/search?name=${query}`);
      console.log(response.data,"====================")
      // Since the response contains both a message and the actual data
      if (Array.isArray(response.data)) {
        setResults(response.data);  // Set the actual character data
      } else {
        setError("No results found.");
      }
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-4">
      <input
        type="text"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search characters..."
        className="border p-2 rounded"
      />
      <button onClick={searchCharacters} className="bg-blue-500 text-white p-2 rounded ml-2">
        Search
      </button>
      {loading && <p>Loading...</p>}
      {error && <p className="text-red-500">{error}</p>}
      <div className="grid grid-cols-3 gap-4 mt-4">
        {results.map((char) => (
          <div key={char._id} className="border p-4 rounded shadow">
            <img src={char.image} alt={char.name} className="w-full h-32 object-cover rounded" />
            <h3 className="text-xl">{char.name}</h3>
            <p>Status: {char.status}</p>
            <p>Species: {char.species}</p>
            <p>Gender: {char.gender}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default CharacterSearch;
