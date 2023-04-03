# modify environment variables in Windows
import os

def export_with_IDA(DIAPHORA_EXPORT_FILE):
  env = os.environ
  env["DIAPHORA_AUTO"] = "1"
  env["DIAPHORA_EXPORT_FILE"] = DIAPHORA_EXPORT_FILE
  os.system('set DIAPHORA_AUTO "{}"'.format(env["DIAPHORA_AUTO"]))
  os.system('set DIAPHORA_EXPORT_FILE "{}"'.format(env["DIAPHORA_EXPORT_FILE"]))

def diff_in_batch_mode(DIAPHORA_DIR):
  env = os.environ
  env["DIAPHORA_AUTO"] = "1"
  env["DIAPHORA_AUTO_DIFF"] = "1"
  env["DIAPHORA_DIFF_OUT"] = DIAPHORA_DIR + "/result.db"
  env["DIAPHORA_DB1"] = DIAPHORA_DIR + "/ori.db"
  env["DIAPHORA_DB2"] = DIAPHORA_DIR + "/modi.db"
  os.system('set DIAPHORA_AUTO "{}"'.format(env["DIAPHORA_AUTO"]))
  os.system('set {} "{}"'.format("DIAPHORA_AUTO_DIFF", env["DIAPHORA_AUTO_DIFF"]))
  os.system('set {} "{}"'.format("DIAPHORA_DIFF_OUT", env["DIAPHORA_DIFF_OUT"]))
  os.system('set {} "{}"'.format("DIAPHORA_DB1", env["DIAPHORA_DB1"]))
  os.system('set {} "{}"'.format("DIAPHORA_DB2", env["DIAPHORA_DB2"]))

def set_timeout_limit(time):
  env = os.environ
  env["DIAPHORA_TIMEOUT_LIMIT"] = str(time)
  os.system('set {} "{}"'.format("DIAPHORA_TIMEOUT_LIMIT", env["DIAPHORA_TIMEOUT_LIMIT"]))

def reset_variable():
  env = os.environ
  del env["DIAPHORA_AUTO"]
  del env["DIAPHORA_AUTO_DIFF"]
  del env["DIAPHORA_DIFF_OUT"]
  del env["DIAPHORA_DB1"]
  del env["DIAPHORA_DB2"]

if __name__ == "__main__":
  set_timeout_limit(120)