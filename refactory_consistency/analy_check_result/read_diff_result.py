import sqlite3
import os


def read_results(db_path, sim_throushold):
    # Connect to the database file
    if not os.path.exists(db_path):
        return
    conn = sqlite3.connect(db_path)

    # Create a cursor object
    cursor = conn.cursor()

    # Execute an SQL statement to get the list of tables
    cursor.execute("SELECT name FROM sqlite_master WHERE type='table';")

    # Fetch all the table names
    tables = cursor.fetchall()

    if len(tables) == 0:
        return False

    # Loop through the table names and retrieve their content
    for table in tables:
        # Execute an SQL statement to retrieve the column names from the current table
        cursor.execute("PRAGMA table_info({})".format(table[0]))

        # Fetch all the column names and print them
        if table[0] == "config":
            continue

        # Execute an SQL statement to retrieve all rows from the current table
        cursor.execute("SELECT * FROM {}".format(table[0]))
        rows = cursor.fetchall()
        if table[0] == "results":
            # Fetch all the rows and print them
            for row in rows:
                if row[0] == "best":
                    continue
                elif row[0] == "partial":
                    if float(row[6]) >= sim_throushold:
                        continue
                    else:
                        cursor.close()
                        conn.close()
                        return False
                else:
                    raise RuntimeError("Unknown results type")
        elif table[0] == "unmatched":
            for row in rows:
                if row[0] == "primary":
                    continue
                elif row[0] == "secondary":
                    cursor.close()
                    conn.close()
                    return False
                else:
                    raise RuntimeError("Unknown unmatched type")
        else:
            raise RuntimeError("Unknown table")
    cursor.close()
    conn.close()
    return True
    # Close the cursor and the connection

if __name__ == "__main__":
  bytecode_dir  = ""
  os.chdir(bytecode_dir )
  proxy_addr_list = os.listdir(os.getcwd())
  for simi in [79]:
    num = 0
    for addr in proxy_addr_list:
        needCount = True
        proxy_dir = bytecode_dir  + "\\" + addr
        impl_list = os.listdir(proxy_dir)
        for impl in impl_list:
            impl_dir = bytecode_dir  + "\\" + addr + "\\" + impl + "\\" + "result.db"
            if not read_results(impl_dir, simi / 100):
                needCount = False
        if needCount:
            num += 1
            print("refactor consistency: ", addr)
        else:
            print("break refactor consistency: ", addr)
