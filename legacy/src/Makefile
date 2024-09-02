.PHONY: init run test

init:
	PYTHONPATH=backend python3 -c "from database import init_db; init_db()"

run:
	PYTHONPATH=backend python3 -c "from database import check_db_exists; check_db_exists()" || make init
	PYTHONPATH=backend python3 backend/app.py

test:
	PYTHONPATH=backend python3 backend/app_tests.py