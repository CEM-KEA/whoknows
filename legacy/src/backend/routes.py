from flask import request, session, url_for, redirect, render_template, flash, jsonify, g
from database import query_db, get_user_id
from security import hash_password, verify_password

def register_routes(app):
    """
    Register all routes with the Flask application.

    Args:
        app (Flask): The Flask application instance.
    """

    @app.route('/')
    def search():
        """
        Handle the search route.

        Returns:
            str: Rendered HTML template for search results.
        """
        q = request.args.get('q')
        language = request.args.get('language', "en")
        if not q:
            search_results = []
        else:
            search_results = query_db("SELECT * FROM pages WHERE language = ? AND content LIKE ?", (language, f"%{q}%"))
        return render_template('search.html', search_results=search_results, query=q)

    @app.route('/about')
    def about():
        """
        Handle the about route.

        Returns:
            str: Rendered HTML template for the about page.
        """
        return render_template('about.html')

    @app.route('/login')
    def login():
        """
        Handle the login route.

        Returns:
            str: Rendered HTML template for the login page.
        """
        if g.user:
            return redirect(url_for('search'))
        return render_template('login.html')

    @app.route('/register')
    def register():
        """
        Handle the register route.

        Returns:
            str: Rendered HTML template for the registration page.
        """
        if g.user:
            return redirect(url_for('search'))
        return render_template('register.html')

    @app.route('/logout')
    def logout():
        """
        Handle the logout route.

        Returns:
            str: Redirect to the search page.
        """
        flash('You were logged out')
        session.pop('user_id', None)
        return redirect(url_for('search'))

    @app.route('/api/search')
    def api_search():
        """
        Handle the API search route.

        Returns:
            Response: JSON response with search results.
        """
        q = request.args.get('q')
        language = request.args.get('language', "en")
        if not q:
            search_results = []
        else:
            search_results = query_db("SELECT * FROM pages WHERE language = ? AND content LIKE ?", (language, f"%{q}%"))
        return jsonify(search_results=search_results)

    @app.route('/api/login', methods=['POST'])
    def api_login():
        """
        Handle the API login route.

        Returns:
            str: Rendered HTML template for the login page or redirect to the search page.
        """
        error = None
        user = query_db("SELECT * FROM users WHERE username = ?", (request.form['username'],), one=True)
        if user is None:
            error = 'Invalid username'
        elif not verify_password(user['password'], request.form['password']):
            error = 'Invalid password'
        else:
            flash('You were logged in')
            session['user_id'] = user['id']
            return redirect(url_for('search'))
        return render_template('login.html', error=error)

    @app.route('/api/register', methods=['POST'])
    def api_register():
        """
        Handle the API register route.

        Returns:
            str: Rendered HTML template for the registration page or redirect to the login page.
        """
        if g.user:
            return redirect(url_for('search'))
        error = None
        if not request.form['username']:
            error = 'You have to enter a username'
        elif not request.form['email'] or '@' not in request.form['email']:
            error = 'You have to enter a valid email address'
        elif not request.form['password']:
            error = 'You have to enter a password'
        elif request.form['password'] != request.form['password2']:
            error = 'The two passwords do not match'
        elif get_user_id(request.form['username']) is not None:
            error = 'The username is already taken'
        else:
            g.db.execute("INSERT INTO users (username, email, password) values (?, ?, ?)",
                         (request.form['username'], request.form['email'], hash_password(request.form['password'])))
            g.db.commit()
            flash('You were successfully registered and can login now')
            return redirect(url_for('login'))
        return render_template('register.html', error=error)