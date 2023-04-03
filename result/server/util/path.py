def some_path(path, position):
  folders = str.split(path, "/")
  positionvalue = folders[len(folders) + position]
  return positionvalue