from flask import g, session
from database import connect_db, query_db

def before_request():
    """
    This function is responsible for connecting to the database and setting the user.
    """
    g.db = connect_db()
    g.user = None
    if 'user_id' in session:
        g.user = query_db("SELECT * FROM users WHERE id = ?", (session['user_id'],), one=True)

def after_request(response):
    """
    This function is responsible for closing the database connection.
    """
    g.db.close()
    return response

def test_before_request_user_id_in_session_but_user_not_in_db(self):
    """
    Ensure that when user_id is in the session but the user doesn't exist in the database,
    g.user is set to None.
    """
    with self.app.test_request_context():
        with self.app.test_client() as client:
            with client.session_transaction() as sess:
                sess['user_id'] = 999  # Assuming 999 doesn't exist in the database
            before_request()
            self.assertIsNone(g.user)

def test_after_request_without_db_connection(self):
    """
    Ensure the after_request function works without any db connection being set in g.
    """
    with self.app.test_request_context():
        response = self.app.response_class()
        try:
            after_request(response)
            self.assertTrue(True)  # No error should occur
        except AttributeError:
            self.fail("after_request raised an AttributeError unexpectedly!")