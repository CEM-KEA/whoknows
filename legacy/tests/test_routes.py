import os
import unittest

from app import app
from test_database import init_test_db

class TestAppRoutes(unittest.TestCase):
    """
    Unit tests for the application's routes.
    """

    def setUp(self):
        """
        Set up the test client and initialize the database.
        """
        self.client = app.test_client()
        self.client.testing = True

        # Initialize the test database
        with app.app_context():
            test_db_path = os.path.abspath('legacy/db/test_whoknows.db')
            os.environ['DATABASE_PATH'] = test_db_path
            init_test_db()  # Initialize the test database

    def tearDown(self):
        """
        Clean up after each test.
        """
        if os.path.exists(os.getenv('DATABASE_PATH')):
            os.remove(os.getenv('DATABASE_PATH'))  # Remove the test database after tests

    def test_login_page_loads_correctly(self):
        """
        Ensure the login page loads correctly.
        """
        response = self.client.get('/login')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Log In', response.data)

    def test_search_page_loads_correctly(self):
        """
        Ensure the search page loads correctly.
        """
        response = self.client.get('/?q=test')
        self.assertEqual(response.status_code, 200)
        # If no specific result is found, you can check for the page structure or a default message
        self.assertIn(b'<input id="search-input"', response.data)  # Checking the search input field exists

    def test_invalid_search_page_loads_correctly(self):
        """
        Ensure the search page loads correctly for an invalid search query.
        """
        response = self.client.get('/?q=no results')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'No search results found.', response.data)

    def test_about_page_loads_correctly(self):
        """
        Ensure the about page loads correctly.
        """
        response = self.client.get('/about')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'About', response.data)

    def test_register_page_loads_correctly(self):
        """
        Ensure the register page loads correctly.
        """
        response = self.client.get('/register')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Register', response.data)

    def test_logout_redirects_to_search(self):
        """
        Ensure logging out redirects to the search page.
        """
        with self.client as c:
            with c.session_transaction() as sess:
                sess['user_id'] = 1
            response = c.get('/logout', follow_redirects=True)
            self.assertEqual(response.status_code, 200)
            self.assertIn(b'You were logged out', response.data)

    def test_api_search_returns_results(self):
        """
        Ensure the API search returns results.
        """
        response = self.client.get('/api/search?q=test')
        self.assertEqual(response.status_code, 200)
        self.assertIn('search_results', response.json)

    def test_api_login_successful(self):
        """
        Ensure the API login is successful with valid credentials.
        """
        with self.client as c:
            # Create the user first
            c.post('/api/register',
                   data=dict(username='testuser', email='test@example.com', password='password', password2='password'))

            response = c.post('/api/login', data=dict(username='testuser', password='password'), follow_redirects=True)
            self.assertEqual(response.status_code, 200)
            self.assertIn(b'You were logged in', response.data)
            with c.session_transaction() as session:
                self.assertIsNotNone(session.get('user_id'))

    def test_api_login_invalid_username(self):
        """
        Ensure the API login fails with an invalid username.
        """
        response = self.client.post('/api/login', data=dict(username='invaliduser', password='password'))
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Invalid username', response.data)

    def test_api_login_invalid_password(self):
        """
        Ensure the API login fails with an invalid password.
        """
        # Ensure user is present in the database
        with self.client as c:
            c.post('/api/register',
                   data=dict(username='testuser', email='test@example.com', password='password', password2='password'))

        response = self.client.post('/api/login', data=dict(username='testuser', password='wrongpassword'))
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Invalid password', response.data)

    def test_api_register_successful(self):
        """
        Ensure the API registration is successful with valid data.
        """
        response = self.client.post('/api/register', data=dict(username='newuser', email='newuser@example.com', password='password', password2='password'), follow_redirects=True)
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'You were successfully registered', response.data)

    def test_api_register_username_taken(self):
        """
        Ensure the API registration fails if the username is already taken.
        """
        # Register the user first
        self.client.post('/api/register',
                         data=dict(username='testuser', email='testuser@example.com', password='password',
                                   password2='password'))

        # Attempt to register again with the same username
        response = self.client.post('/api/register',
                                    data=dict(username='testuser', email='testuser@example.com', password='password',
                                              password2='password'))
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'The username is already taken', response.data)

    def test_api_register_invalid_email(self):
        """
        Ensure the API registration fails with an invalid email address.
        """
        response = self.client.post('/api/register', data=dict(username='newuser', email='invalidemail', password='password', password2='password'))
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'You have to enter a valid email address', response.data)

    def test_api_register_passwords_do_not_match(self):
        """
        Ensure the API registration fails if the passwords do not match.
        """
        response = self.client.post('/api/register', data=dict(username='newuser', email='newuser@example.com', password='password', password2='differentpassword'))
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'The two passwords do not match', response.data)

    def test_invalid_route_returns_404(self):
        """
        Ensure an invalid route returns a 404 status code.
        """
        response = self.client.get('/nonexistent')
        self.assertEqual(response.status_code, 404)
        self.assertIn(b'404', response.data)