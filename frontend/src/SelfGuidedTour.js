import React, { useEffect, useState } from "react";
import axios from "axios";
import "./style/CityGuidePage.css"; // Import the CSS file
import BackToHomeButton from "./BackToHomeButton"; // Import the BackToHomeButton component


function CityGuidePage() {
  const [cityTours, setCityTours] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [connectedToBackend, setConnectedToBackend] = useState(false);
  const [error, setError] = useState(null);

  // Define the function to handle guide selection
  const handleTourSelection = (selectedGuide) => {
    // You can add your logic here to handle the selected guide
    console.log("Selected Guide:", selectedGuide);
  };

  useEffect(() => {
    // Check if the frontend is connected to the backend
    axios
      .get("http://localhost:3001/api/status")
      .then(() => {
        setConnectedToBackend(true);
      })
      .catch((error) => {
        console.error("Error connecting to backend:", error);
      });

    // Fetch city guides from the backend when the component mounts
    axios
      .get("http://localhost:3001/api/citytour")
      .then((response) => {
        if (response.data) {
          setCityTours(response.data);
          setIsLoading(false);
        }
      })
      .catch((error) => {
        setError(error);
        setIsLoading(false);
      });
  }, []);

  if (isLoading) {
    return <p>Loading...</p>;
  }

  if (error) {
    return <p>Error: {error.message}</p>;
  }

  return (
    <div className="container">
      <h1>City Tours</h1>
      <ul>
        {cityTours.map((tour) => (
          <li key={tour.name} className="block tour" onClick={() => handleTourSelection(tour)}>
            <strong>Name:</strong> {tour.name} <br />
            <strong>Contact:</strong> {tour.contact} <br />
            <strong>Price:</strong> {tour.price} <br />
            <strong>Personal Data:</strong> {tour.personal_data} <br />
            {/* Add a button or link to select the guide */}
            <button className="select-button">Select Tour</button>
          </li>
        ))}
      </ul>
    <BackToHomeButton />  
    </div>
  );
}

export default CityGuidePage;
