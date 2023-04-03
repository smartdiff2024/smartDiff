import os
import csv


diffproxyaddress = []
resultDir = "./dir"
os.mkdir(resultDir)
with open("./order", "r", encoding="utf-8") as csvFile:
  reader = csv.reader(csvFile)
  header = next(reader)
  for row in reader:
    if row[0] not in diffproxyaddress:
      diffproxyaddress.append(row[0])
      os.mknod(resultDir + "/" + row[0] + ".csv")
    writeFile = open(resultDir + "/" + row[0] + ".csv", "a", encoding="utf-8")
    writer = csv.writer(writeFile)
    writer.writerow(row)
    writeFile.close()