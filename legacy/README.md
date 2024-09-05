# ¿Who Knows? - Flask Application

A legacy search engine built in 2009, now ported to Python 3 and using more modern Flask versions! This project demonstrates how to modernize and work with older Python applications while maintaining backwards compatibility.

**Note**: This application is still intentionally full of problems and vulnerabilities for educational purposes. **Do not run this in a production environment**.

## Requirements

- Python 3.6+
- SQLite (for local database use)
- Virtual Environment (recommended)

## Installation

It is highly recommended to use a **virtual environment** to manage dependencies:

```bash
python3 -m venv .venv
source .venv/bin/activate  # On Windows, use .venv\Scripts\activate
```

Install the necessary dependencies:
    
```bash
make install
```
This will create the virtual environment and install the dependencies from the requirements.txt file.

## Database Setup
To initialize a new database, run:

```bash
make init
```

This will create the SQLite database (whoknows.db) based on the schema provided in schema.sql.

**Note**: This will delete any existing data in the database.
**Note**: If you make changes to the database schema, you'll need to re-run make init.

## Running the Application
To start the development server on port 8080, run:
    
```bash
make run
```

Alternatively, you can run the application directly using:

```bash
python3 legacy/src/backend/app.py
```

**Note**: Ensure the virtual environment is activated before running the application.

## Tests
To run the tests, use:

```bash
make test
```
Or run the tests directly using:

```bash
python3 -m unittest discover -s tests
```

## Development Notes
- This application was originally built using Python 2.7 and Flask 0.5, but it has been ported to work with Python 3 and the latest Flask versions.
- Security Warning: This application contains intentional security flaws for educational purposes. Avoid using this in any real-world environment.
- The database is set up using SQLite. You can find the schema in legacy/db/schema.sql.
- The project structure has been refactored to separate the application code from the configuration and testing components.

## Windows Users
Windows does not natively support Make. If you are using Windows, you can manually run the corresponding commands. For example, instead of make init, you would run:

```bash
source .venv/bin/activate
python3 -c "from database import init_db; init_db()"
```
## License
MIT License

Copyright (c) 2024 who-knows-inc

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.