import csv

dtrace = []

def readtrace(csvPath):
  with open(csvPath, "r", encoding="utf-8") as csvFile:
    reader = csv.reader(csvFile)
    for row in reader:
      dtrace.append(row)

def errorFirstImpe(proxyaddress):
  impls = []
  for row in dtrace:
    if str.lower(row[5]) == str.lower(proxyaddress):
      for index in range(8, len(row), 4):
        impl_low = str.lower(row[index])
        if impl_low not in impls:
          impls.append(impl_low)
        else:
          print("Two Times Error: ", row)
          return True
  return False

if __name__ == "__main__":
  readtrace("./trace.csv")
      