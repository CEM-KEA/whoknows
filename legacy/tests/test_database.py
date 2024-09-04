import unittest
import os
import sys

# Add the directory containing the `database` module to the Python path
sys.path.append(os.path.join(os.path.dirname(__file__), '../../src/backend'))
from flask import g
from database import connect_db, init_db, get_user_id, query_db, check_db_exists
from app import app

def init_test_db():
    """
    Initializes the test database with the schema.
    """
    with open(os.path.join(os.path.dirname(__file__), '../db/schema.sql')) as schema_file:
        conn = connect_db(init_mode=True)  # Create connection to the test database
        conn.executescript(schema_file.read())  # Initialize the schema from schema.sql
        conn.commit()  # Commit the changes
        conn.close()  # Close the connection


class TestDatabase(unittest.TestCase):
    """
    Unit tests for the database functions.
    """

    def setUp(self):
        """
        Set up the test environment before each test.
        """
        test_db_path = os.path.abspath('legacy/db/test_whoknows.db')
        os.environ['DATABASE_PATH'] = test_db_path
        os.makedirs(os.path.dirname(test_db_path), exist_ok=True)
        with app.app_context():
            init_test_db()

    def tearDown(self):
        """
        Clean up the test environment after each test.
        """
        if os.path.exists(os.getenv('DATABASE_PATH')):
            os.remove(os.getenv('DATABASE_PATH'))

    def test_get_user_id(self):
        """
        Test that the correct user ID is retrieved.
        """
        with app.app_context():
            conn = connect_db()
            conn.execute("INSERT INTO users (username, email, password) values (?, ?, ?)",
                         ('testuser', 'test@example.com', 'password'))
            conn.commit()
            user_id = get_user_id('testuser', conn)
            self.assertIsNotNone(user_id)
            self.assertIsInstance(user_id, int)
            conn.close()

    def test_get_user_id_nonexistent_user(self):
        """
        Test that None is returned for a nonexistent user.
        """
        with app.app_context():
            conn = connect_db()
            user_id = get_user_id('nonexistent', conn)
            self.assertIsNone(user_id)
            conn.close()

    def test_init_db(self):
        """
        Test that init_db initializes the database schema correctly.
        """
        with app.app_context():
            conn = connect_db()
            # Check if the 'users' table exists
            result = conn.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='users';").fetchone()
            self.assertIsNotNone(result)
            conn.close()

    def test_query_db(self):
        """
        Test the query_db function for correct behavior.
        """
        with app.app_context():
            conn = connect_db()
            conn.execute("INSERT INTO users (username, email, password) values (?, ?, ?)",
                         ('testuser', 'test@example.com', 'password'))
            conn.commit()
            # Query the user back using query_db
            g.db = conn  # Manually set g.db for the query_db function to use
            result = query_db("SELECT * FROM users WHERE username = ?", ['testuser'], one=True)
            self.assertIsNotNone(result)
            self.assertEqual(result['username'], 'testuser')
            conn.close()

    def test_check_db_exists(self):
        """
        Test that check_db_exists correctly identifies the database existence.
        """
        self.assertTrue(check_db_exists())
        os.remove(os.getenv('DATABASE_PATH'))
        with self.assertRaises(RuntimeError):
            check_db_exists()

    def test_connect_db_error(self):
        """
        Test that connect_db raises an exception when the database does not exist.
        """
        os.remove(os.getenv('DATABASE_PATH'))
        with self.assertRaises(RuntimeError):
            connect_db()

if __name__ == '__main__':
    unittest.main()