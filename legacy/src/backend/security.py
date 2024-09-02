import bcrypt

def hash_password(password):
    """
    Hashes a password using bcrypt.
    """
    password_bytes = password.encode('utf-8')
    salt = bcrypt.gensalt()
    return bcrypt.hashpw(password_bytes, salt)

def verify_password(hashed_password, password):
    """
    Verifies a password against a hashed password using bcrypt
    """
    hashed_password_bytes = hashed_password.encode('utf-8')
    password_bytes = password.encode('utf-8')
    return bcrypt.checkpw(password_bytes, hashed_password_bytes)