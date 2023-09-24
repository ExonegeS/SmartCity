import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login, logout } from './auth';
import "./style/LoginPage.css"; // Import the CSS file
import BackToHomeButton from "./BackToHomeButton"; // Import the BackToHomeButton component


function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate(); // React Router hook for navigation

  const handleLogin = async (e) => {
    e.preventDefault();

    try {
      // Call the login function with the entered email and password
      await login(email, password);

      // Redirect to the desired page after successful login
      navigate('/'); // Redirect to the admin page
    } catch (error) {
      navigate('/'); // Redirect to the main page
      // Handle login error (e.g., show an error message)
    }
  };

  return (
    <div>
      <h2>Login</h2>
      <form onSubmit={handleLogin}>
        <div>
          <label>Email:</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Password:</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit">Login</button>
      </form>
      <BackToHomeButton/>
    </div>
  );
}

export default LoginPage;
