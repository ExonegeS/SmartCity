import React, { useEffect } from "react";
import { Link } from "react-router-dom";
import { isAuthenticated, logout } from './auth'; // Import isAuthenticated and logout functions
import "./style/MainPage.css";

function MainPage() {
  // Check if the user is authenticated
  const userAuthenticated = isAuthenticated();

  // Define the handleLogout function
  const handleLogout = () => {
    // Call the logout function to log the user out
    logout();
    window.location.reload(false);
  };

  useEffect(() => {
    // Scroll to a specific element (e.g., the first Link block) when the component mounts
    const element = document.getElementById("initial-scroll-target");
    if (element) {
      element.scrollIntoView({ behavior: "smooth", block: "center" });
    }
  }, []);

  return (
    <div className="page-container">
      <div className="content-container">
        <h1>Welcome to My Smart City</h1>
        <p>Explore the city's features:</p>
        <Link to="/cityguide" className="block">
          <div className="block-content">
            <h2>City Guide</h2>
            <p>Description for this block...</p>
          </div>
          <div className="block-image">
            <img src="images/city_guide.svg" alt="City Guide" />
          </div>
        </Link>
        {/* Conditionally render the link based on authentication status */}
          {userAuthenticated ? (
            <Link to="/admin/cityguide" className="block admin">
            <h2>Edit City Guides</h2></Link>
          ) : null}

        <Link to="/selfguidedtour" className="block">
          <div className="block-content">
            <h2>Self-guided city tour</h2>
            <p>Description for this block...</p>
          </div>
          <div className="block-image">
            <img src="images/self_guide.svg" alt="Self-guided Tour" />
          </div>
        </Link>

        <Link to="/photooftheday" className="block">
          <div className="block-content">
            <h2>Photo of the day</h2>
            <p>Description for this block...</p>
          </div>
          <div className="block-image">
            <img src="images/photo.svg" alt="Photo of the Day" />
          </div>
        </Link>

        <Link to="/localdishes" className="block">
          <div className="block-content">
            <h2>Challenge to try local dishes of the city</h2>
            <p>Description for this block...</p>
          </div>
          <div className="block-image">
            <img src="images/food.svg" alt="Local Dishes" />
          </div>
        </Link>

        <Link to="/historicalpart" className="block">
          <div className="block-content">
            <h2>Historical part</h2>
            <p>Description for this block...</p>
          </div>
          <div className="block-image">
            <img src="images/clock.svg" alt="Historical Part" />
          </div>
        </Link>
        
        <Link to="/blank" className="block">
          <div className="block-content">
            <h2>Blank part</h2>
            <p>Description for this blank block...</p>
          </div>
          <div className="block-image">
            <img src="image6.jpg" alt="Blank" />
          </div>
        </Link>
        {/* Conditionally render the login/logout link based on authentication status */}
        {userAuthenticated ? (
              <button onClick={handleLogout}>Logout</button>
            ):<></>}
      </div>
    </div>
  );
}

export default MainPage;
