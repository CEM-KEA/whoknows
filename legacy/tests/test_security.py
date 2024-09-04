import unittest
from security import hash_password, verify_password

class TestSecurity(unittest.TestCase):
    """
    Unit tests for the security functions `hash_password` and `verify_password`.
    """

    def test_hash_password(self):
        """
        Test that hashing a password returns a hash different from the original password
        and that the hash can be verified against the original password.
        """
        password = "testpassword"
        hashed_pw = hash_password(password)
        self.assertNotEqual(password, hashed_pw)
        self.assertTrue(verify_password(hashed_pw.decode('utf-8'), password))

    def test_verify_password(self):
        """
        Test that verifying a password works as expected for both correct and incorrect passwords.
        """
        password = "testpassword"
        hashed_pw = hash_password(password)
        self.assertTrue(verify_password(hashed_pw.decode('utf-8'), password))
        self.assertFalse(verify_password(hashed_pw.decode('utf-8'), "wrongpassword"))

    def test_empty_password(self):
        """
        Test that hashing and verifying an empty password works correctly.
        """
        password = ""
        hashed_pw = hash_password(password)
        self.assertNotEqual(password, hashed_pw)
        self.assertTrue(verify_password(hashed_pw.decode('utf-8'), password))

    def test_long_password(self):
        """
        Test that hashing and verifying a very long password works correctly.
        """
        password = "a" * 500
        hashed_pw = hash_password(password)
        self.assertNotEqual(password, hashed_pw)
        self.assertTrue(verify_password(hashed_pw.decode('utf-8'), password))

    def test_unicode_password(self):
        """
        Test that hashing and verifying a password with unicode characters works correctly.
        """
        password = "pässwørd"
        hashed_pw = hash_password(password)
        self.assertNotEqual(password, hashed_pw)
        self.assertTrue(verify_password(hashed_pw.decode('utf-8'), password))

    def test_invalid_hashed_password(self):
        """
        Test that verifying a password against an invalid hash returns False.
        """
        password = "testpassword"
        invalid_hashed_pw = "notahash"
        self.assertFalse(verify_password(invalid_hashed_pw, password))

    def test_special_character_password(self):
        """
        Test that hashing and verifying a password with special characters works correctly.
        """
        password = "!@#$%^&*()_+-=~`"
        hashed_pw = hash_password(password)
        self.assertNotEqual(password, hashed_pw)
        self.assertTrue(verify_password(hashed_pw.decode('utf-8'), password))

    def test_none_password(self):
        """
        Test that None input is handled properly for both hash and verify functions.
        """
        with self.assertRaises(ValueError):
            hash_password(None)
        self.assertFalse(verify_password(None, "testpassword"))

    def test_hash_uniqueness(self):
        """
        Test that hashing the same password multiple times produces unique hashes.
        """
        password = "samepassword"
        hashed_pw1 = hash_password(password)
        hashed_pw2 = hash_password(password)
        self.assertNotEqual(hashed_pw1, hashed_pw2)

if __name__ == '__main__':
    unittest.main()