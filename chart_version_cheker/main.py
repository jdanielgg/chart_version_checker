""" Usage:
  chart_check_version <appName> <chart> <chart>
  check_index_version (-h | --help)
  chart_check_version --version

Options:
  -h --help     Show this screen.
  --version     Show version.
"""
from docopt import docopt
import yaml
import semver
import sys
from functools import reduce


def main(appName, chartFile, indexFile):
    chart = parse(read(chartFile)['version'])

    index = read(indexFile)
    latest = getLatest(appName, IndexInfo=index)

    if chart <= latest:
        # we took the exit code forme https://tldp.org/LDP/abs/html/exitcodes.html
        print("The chart version should be greather than the lastes published verion", file=sys.stderr)
        sys.exit(126)


def getLatest(appName, IndexInfo):

    if 'entries' not in IndexInfo.keys():
        print("Malformed index yaml", file=sys.stderr)
        sys.exit(1)

    if appName not in IndexInfo['entries'].keys():
        print(f"App {appName} not found on index", file=sys.stderr)
        sys.exit(1)

    if 'version' not in IndexInfo['entries'][appName].keys():
        print("Malformed index yaml (version not found)", file=sys.stderr)
        sys.exit(1)

    versions = [ parse(x['version']) for x in IndexInfo['entries'][appName]]
    return reduce(max, versions)


def read(name):
    with open(name, 'r') as f:
        text = f.read()

    return yaml.safe_load(text)


def parse(aVersion):
    return semver.VersionInfo.parse(aVersion)


if __name__ == "__main__":
    args = docopt(__doc__, version='1.0.0')
    main(
        appName=args['<appName>'],
        chartFile=args['<chart>'],
        indexFile=args['<index>'],
    )
