from flask import g

from handlers import before_request, after_request

def before_request_sets_db_connection(self):
    """
    Ensure the before_request function sets the database connection.
    """
    with self.app.test_request_context():
        before_request()
        self.assertIsNotNone(g.db)

def before_request_sets_user_when_logged_in(self):
    """
    Ensure the before_request function sets the user when logged in.
    """
    with self.app.test_request_context():
        with self.app.test_client() as client:
            with client.session_transaction() as sess:
                sess['user_id'] = 1
            before_request()
            self.assertIsNotNone(g.user)

def before_request_does_not_set_user_when_not_logged_in(self):
    """
    Ensure the before_request function does not set the user when not logged in.
    """
    with self.app.test_request_context():
        before_request()
        self.assertIsNone(g.user)

def after_request_closes_db_connection(self):
    """
    Ensure the after_request function closes the database connection.
    """
    with self.app.test_request_context():
        before_request()
        response = self.app.response_class()
        after_request(response)
        self.assertTrue(g.db.closed)