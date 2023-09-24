import axios from 'axios';

// Example login function
const login = async (email, password) => {
  try {
    const response = await axios.post("http://localhost:3001/api/login", { email, password });

    const token = response.data.token;
    // Store the token in a secure manner (e.g., in a cookie or local storage)
    localStorage.setItem('token', token); // Store the token in local storage

    // Use the token for future requests
    axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;
  } catch (error) {
    console.error("Login error");
    throw error; // Rethrow the error for handling in the LoginPage component
  }
};

// Example logout function
const logout = () => {
  // Clear the token from storage
  localStorage.removeItem('token');

  // Remove the token from axios headers
  delete axios.defaults.headers.common["Authorization"];
};

// Example function to check if a user is authenticated
const isAuthenticated = () => {
  // Check if the token exists in local storage or implement your own logic
  const token = localStorage.getItem('token');
  return !!token;
};

export { login, logout, isAuthenticated };
