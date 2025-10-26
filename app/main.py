import os
import datetime

def main():
    log_file = os.path.expanduser('~/shuttlexpress.log')
    with open(log_file, 'a') as f:
        f.write(f'ShuttleXpress connected at {datetime.datetime.now()}\n')

if __name__ == '__main__':
    main()
