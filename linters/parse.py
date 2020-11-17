import sqlparse
import sys


def main(argv):
  if len(argv)<=1:
    print("请输入格式化文件地址")
    return
  for name in argv[1:]:
    f = open(name)
    lines = f.read()
    f.close()
    f = open(name,"w")
    f.write(sqlparse.format(lines, reindent=True, keyword_case='upper'))

if __name__ == "__main__":
   main(sys.argv)
