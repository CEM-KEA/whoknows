# Whoknows

This is a legacy search engine that has been ported from Python 2 to Python 3.

## Dependency graph
Can be viewed with a dependency graph visual editor, ex. http://magjac.com/graphviz-visual-editor/
```
strict digraph {
    app -> Flask
    app -> handlers
    app -> routes
    app -> dotenv
    app -> database

    database -> os
    database -> sqlite3
    database -> contextlib
    database -> dotenv
    database -> flask
    database -> app

    handlers -> flask
    handlers -> database

    routes -> flask
    routes -> database
    routes -> security

    security -> bcrypt

    test_routes -> unittest
    test_routes -> os
    test_routes -> app
    test_routes -> test_database

    test_database -> unittest
    test_database -> os
    test_database -> sys
    test_database -> flask
    test_database -> database
    test_database -> app
}
```
