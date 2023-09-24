import React from "react";
import { Link } from "react-router-dom";
import "./style/BackToHomeButton.css"; // Import the CSS file


function BackToHomeButton() {
  return (
    <Link to="/" className="back-to-home-button">
      Back to Home
    </Link>
  );
}

export default BackToHomeButton;
