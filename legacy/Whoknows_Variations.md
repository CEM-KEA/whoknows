# Whoknows Variations

This is the Whoknows variations repository. It is not meant for production as it contains several security vulnerabilities and problematic parts on purpose. 

## How to get started

Each branch is a tutorial in a different topic based on the same Flask application as in the `main` branch. 

One way to follow along is by:

1. Forking the repository to your own account.

2. Cloning the repository to your local machine.

3. Checking out the branch you are interested in (e.g. `git checkout <branch_name>`).

4. Following the instructions in the README of the branch.

5. You can now push changes to your own repository. 

## Pull requests

If you have any suggestions or improvements to the tutorials, feel free to open a pull request.

## Dependency graph
Can be viewed with a dependency graph visual editor, ex. http://magjac.com/graphviz-visual-editor/
```
strict digraph {
    server -> flask
    server -> database
    server -> sqlite3
    server -> hashlib
    server -> os
    server -> sys
    flask -> templates
    test -> server
    templates -> css
    templates -> jinja
}
```
