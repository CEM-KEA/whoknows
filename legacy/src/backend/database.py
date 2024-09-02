import os
import sqlite3
import sys
from contextlib import closing
from flask import g
from config import DATABASE_PATH

def connect_db(init_mode=False):
    """
    This function is responsible for connecting to the database.
    """
    if not init_mode:
        check_db_exists()
    return sqlite3.connect(DATABASE_PATH)

def check_db_exists():
    """
    This function is responsible for checking if the database exists.
    """
    db_exists = os.path.exists(DATABASE_PATH)
    if not db_exists:
        print("Database not found")
        sys.exit(1)
    else:
        return db_exists

def init_db():
    """
    This function is responsible for initializing the database.
    """
    from app import app  # Import here to avoid circular dependency
    with closing(connect_db(init_mode=True)) as db:
        with app.open_resource('../schema.sql') as f:
            db.cursor().executescript(f.read().decode())
        db.commit()
        print("Initialized the database: " + str(DATABASE_PATH))

def query_db(query, args=(), one=False):
    """
    This function is responsible for querying the database.
    """
    cur = g.db.execute(query, args)
    rv = [dict((cur.description[idx][0], value)
               for idx, value in enumerate(row)) for row in cur.fetchall()]
    return (rv[0] if rv else None) if one else rv

def get_user_id(username):
    """
    This function is responsible for getting the user ID.
    """
    rv = g.db.execute("SELECT id FROM users WHERE username = ?", (username,)).fetchone()
    return rv[0] if rv else None