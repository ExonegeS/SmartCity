import React, { useState, useEffect } from "react";
import axios from "axios";
import { isAuthenticated } from './auth'; // Import the isAuthenticated function
import BackToHomeButton from "./BackToHomeButton"; // Import the BackToHomeButton component


function AdminCityGuidePage() {
  const [newGuideData, setNewGuideData] = useState({
    name: "",
    contact: "",
    price: 0,
    personal_data: "",
  });

  const [cityGuides, setCityGuides] = useState([]);

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setNewGuideData({ ...newGuideData, [name]: value });
  };

  const handleAddGuide = () => {
    // Check if the user is authenticated before making the request
    if (isAuthenticated()) {
      // Convert the price to a number
      const priceAsNumber = Number(newGuideData.price);
  
      axios
        .post("http://localhost:3001/api/cityguide", {
          ...newGuideData,
          price: priceAsNumber, // Send the price as a number
        })
        .then((response) => {
          console.log("New guide added:", response.data);
          // Optionally, you can reset the form fields after adding a guide
          setNewGuideData({
            name: "",
            contact: "",
            price: 0,
            personal_data: "",
          });
          // Refresh the list of city guides after adding
          fetchCityGuides();
        })
        .catch((error) => {
          console.error("Error adding new guide:", error);
        });
    } else {
      // Handle the case where the user is not authenticated (e.g., show a message or redirect to the login page)
      console.error("User is not authenticated. Please log in.");
    }
  };

  const handleDeleteGuide = (name) => {
    // Check if the user is authenticated before making the request
    if (isAuthenticated()) {
      axios
        .delete(`http://localhost:3001/api/cityguide/${name}`)
        .then((response) => {
          console.log("City guide deleted:", response.data);
          // Refresh the list of city guides after deleting
          fetchCityGuides();
        })
        .catch((error) => {
          console.error("Error deleting city guide:", error);
        });
    } else {
      // Handle the case where the user is not authenticated (e.g., show a message or redirect to the login page)
      console.error("User is not authenticated. Please log in.");
    }
  };

  useEffect(() => {
    // Fetch the list of city guides when the component mounts
    fetchCityGuides();
  }, []);

  const fetchCityGuides = () => {
    axios
      .get("http://localhost:3001/api/cityguide")
      .then((response) => {
        setCityGuides(response.data);
      })
      .catch((error) => {
        console.error("Error fetching city guides:", error);
      });
  };

  return (
    <div>
      <h2>Admin: Add New City Guide</h2>
      <form>
        {/* Form for adding a new city guide */}
        <div>
          <label>Name:</label>
          <input
            type="text"
            name="name"
            value={newGuideData.name}
            onChange={handleInputChange}
          />
        </div>
        <div>
          <label>Contact:</label>
          <input
            type="text"
            name="contact"
            value={newGuideData.contact}
            onChange={handleInputChange}
          />
        </div>
        <div>
          <label>Price:</label>
          <input
            type="number"
            name="price"
            value={newGuideData.price}
            onChange={handleInputChange}
          />
        </div>
        <div>
          <label>Personal Data:</label>
          <input
            type="text"
            name="personal_data"
            value={newGuideData.personal_data}
            onChange={handleInputChange}
          />
        </div>
        <button type="button" onClick={handleAddGuide}>
          Add Guide
        </button>
      </form>

      {/* Display a list of existing city guides */}
      <h2>Existing City Guides</h2>
      <ul>
        {cityGuides.map((guide) => (
          <li key={guide.name}>
            {guide.name}{" "}
            <button onClick={() => handleDeleteGuide(guide.name)}>Delete</button>
          </li>
        ))}
      </ul>
      <BackToHomeButton/>
    </div>
  );
}

export default AdminCityGuidePage;
