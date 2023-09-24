import React from "react";
import { Link } from "react-router-dom";
import "./style/NotFoundPage.css";

function NotFoundPage() {
  return (
    <div>
      <div className="pattern-container"></div>

      <div className="not-found-container">
        <h2 className="not-found-title">404 - Page Not Found</h2>
        <p className="not-found-text">The page you are looking for does not exist.</p>
        <Link to="/" className="back-button">
          Back to Home
        </Link>
      </div>
    </div>
  );
}

export default NotFoundPage;
