import bcrypt


def hash_password(password):
    """
    Hashes a password using bcrypt.
    """
    if password is None:
        raise ValueError("Password cannot be None")

    password_bytes = password.encode('utf-8')
    salt = bcrypt.gensalt()
    return bcrypt.hashpw(password_bytes, salt)


def verify_password(hashed_password, plain_password):
    """
    Verifies a password against a hashed password using bcrypt
    """
    try:
        # Ensure the hashed password is in bytes
        if isinstance(hashed_password, bytes):
            hashed_password_bytes = hashed_password
        else:
            hashed_password_bytes = hashed_password.encode('utf-8')

        return bcrypt.checkpw(plain_password.encode('utf-8'), hashed_password_bytes)
    except (ValueError, AttributeError):
        # Catch errors such as invalid salt or encoding issues
        return False
