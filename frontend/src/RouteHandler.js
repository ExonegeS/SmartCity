import React, { useContext } from "react";
import { BrowserRouter as Router, Routes, Route, Navigate, useHistory } from "react-router-dom";
import MainPage from "./MainPage";
import LoginPage from "./LoginPage";
import CityGuidePage from "./CityGuidePage";
import AdminCityGuidePage from "./AdminCityGuidePage";
import SelfGuidedTour from "./SelfGuidedTour";
import NotFoundPage from "./NotFoundPage"; // Import the NotFoundPage component
import { isAuthenticated } from "./auth"; // Import isAuthenticated from auth.js

function ProtectedRoute({ element, ...rest }) {
  const isUserAuthenticated = isAuthenticated();

  return isUserAuthenticated ? element : <Navigate to="/login" />;
}

function RouteHandler() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<MainPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/cityguide" element={<CityGuidePage />} />
        {/* Use ProtectedRoute for admin pages */}
        <Route path="/admin/cityguide" element={<ProtectedRoute element={<AdminCityGuidePage />} />} />

        <Route path="/selfguidedtour" element={<SelfGuidedTour />} />
        <Route path="*" element={<NotFoundPage />} /> {/* Use the NotFoundPage component */}
      </Routes>
    </Router>
  );
}

export default RouteHandler;
