import argparse
import contextlib
import os
import random
import shutil
import string
import tempfile
import time
import zipfile
from pathlib import Path

import requests
import tqdm
from packaging import version

BACKUP = Path('.', '.backup')


@contextlib.contextmanager
def temporary_directory(dir: str = None, delete: bool = True):
    def random_string() -> str:
        return ''.join(random.choices(string.ascii_letters + string.digits, k=8))

    dir_temp = Path(dir or tempfile.gettempdir()).joinpath(random_string())

    while dir_temp.exists():
        dir_temp = dir_temp.parent.joinpath(random_string())
    os.mkdir(dir_temp, 0o777)

    try:
        yield dir_temp
    finally:
        if delete:
            shutil.rmtree(dir_temp, True)


def backup():
    if BACKUP.exists():
        shutil.rmtree(BACKUP)
    os.mkdir(BACKUP)

    for filename in ('driver-box.exe', 'bin', 'conf'):
        if not (path := Path(filename)).exists():
            continue
        shutil.copy(path, BACKUP.joinpath(filename))


def cleanup(restore: bool):
    if not restore:
        print('Removing backups...')
    else:
        print('Restoring states...')
        for filename in ('driver-box.exe', 'bin', 'conf'):
            if (newfile := Path(filename)).exists():
                if newfile.is_dir():
                    shutil.rmtree(newfile, True)
                else:
                    newfile.unlink()

            if not BACKUP.joinpath(filename).exists():
                continue

            BACKUP.joinpath(filename).rename(filename)

    shutil.rmtree(BACKUP, True)


def replace_executable(version: str, binary_type: str, webview: bool):
    filename = f'driver-box.{binary_type}-wv2.zip' if webview else f'driver-box.{binary_type}.zip'
    resp = requests.get('https://github.com/markmybytes/driver-box/releases/download/'
                        f'v{version}/{filename}',
                        stream=True)

    if resp.headers.get('content-type') not in ('application/zip', 'application/octet-stream'):
        raise ValueError('invalid version or binary type')

    with temporary_directory(dir=os.getcwd()) as tmpdir:
        tmpdir = Path(tmpdir)
        fpath = tmpdir.joinpath(filename)

        print(f'Downloading: {filename}')
        with (tqdm.tqdm(total=int(resp.headers['Content-Length']), unit="B", unit_scale=True) as progress,
                open(fpath, 'wb') as f):
            for chunk in resp.iter_content(1024):
                f.write(chunk)
                progress.update(len(chunk))
                progress.display()

        print(f'\nUnpacking...')
        with zipfile.ZipFile(fpath, 'r') as z:
            for archive in tqdm.tqdm(z.filelist, unit='file'):
                z.extract(archive.filename, str(tmpdir))

        print('Updating...')
        paths = ('driver-box.exe', 'bin') if webview else ('driver-box.exe',)
        for path in map(Path, paths):
            if path.exists():
                if path.is_dir():
                    shutil.rmtree(path, True)
                else:
                    path.unlink()

            time.sleep(1)  # add wait time to avoid WinError5

            if tmpdir.joinpath(path).exists():
                tmpdir.joinpath(path).rename(Path('.', path))


def migrate_config(from_: version.Version, to: version.Version):
    if from_.major == to.major:
        return

    if from_.major == 1 and to.major == 2 or to.major >= 4:
        raise NotImplementedError()
    return


if __name__ == '__main__':
    print(r'''
     _      _                     _                                 _       _            
  __| |_ __(_)_   _____ _ __     | |__   _____  __  _   _ _ __   __| | __ _| |_ ___ _ __ 
 / _` | '__| \ \ / / _ \ '__|____| '_ \ / _ \ \/ / | | | | '_ \ / _` |/ _` | __/ _ \ '__|
| (_| | |  | |\ V /  __/ | |_____| |_) | (_) >  <  | |_| | |_) | (_| | (_| | ||  __/ |   
 \__,_|_|  |_| \_/ \___|_|       |_.__/ \___/_/\_\  \__,_| .__/ \__,_|\__,_|\__\___|_|   
                                                         |_|                             
''')

    argparser = argparse.ArgumentParser(description='')
    argparser.add_argument('-d', '--app-directory', type=str,
                           help='Root directory of driver-box')
    argparser.add_argument('-s', '--version-from', type=str,
                           required=True, help='Update from which verion')
    argparser.add_argument('-t', '--version-to', type=str,
                           required=True, help='Update to which version')
    argparser.add_argument('-b', '--binary-type', type=str,
                           required=True, help='Binary target')
    argparser.add_argument('-w', '--webview', action='store_true',
                           help='Download built-in WebView2 verion')

    args = argparser.parse_args()
    if args.app_directory:
        os.chdir(args.app_directory)

    version_from = version.parse(args.version_from)
    version_to = version.parse(args.version_to)

    if version_from.major > version_to.major:
        print('Downgrade is not supported!')
        input('Press any key to exit...')
        exit()

    try:
        print('+', '-'*26, '+')
        print('| {:13s}{:^13s} |'.format('Update From', str(version_from)))
        print('| {:13s}{:^13s} |'.format('Update To', str(version_to)))
        print('| {:13s}{:^13s} |'.format('Binary', args.binary_type))
        print('| {:13s}{:^13s} |'.format(
            'WebView2', 'Yes' if args.webview else 'No'))
        print('+', '-'*26, '+', end='\n\n')

        backup()

        replace_executable(str(version_to), args.binary_type, args.webview)

        migrate_config(version_from, version_to)
    except Exception as e:
        print(f'Error occures: {e}')
        cleanup(True)
    else:
        cleanup(False)

    input('Finished. Press any key to continue...')
