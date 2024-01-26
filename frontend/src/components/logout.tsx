export function Logout() {
  async function logout() {
    try {
      const response = await fetch('http://localhost:9090/account/logout', {
        method: 'POST',
        credentials: 'include',
      });

      const data = await response.json();

      // Check for successful login (adjust based on your backend response structure)
      if (response.ok) {
        // Redirect to the dashboard
        window.location.href = '/';
      } else {
        // Handle login failure (display an error message, etc.)
        console.error('failed to logout:', data.error);
        alert('failed to logout:'+ String(data.error));
      }
    } catch (error) {
      console.error('Error during logging out:', error);
      alert('failed to logout:'+ String(error));
    }
  }
  logout()
}