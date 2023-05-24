import requests
import os
from typing import Optional

def get_parameter(variable: str, default: str) -> str:
    """
    Get parameter from the given environment variable if present, otherwise return the default
    :param str variable: The environment variable to check
    :param str default: The value to return if the environment variable is not provided
    :returns
    """
    env_value = os.environ.get(variable)

    if env_value != None:
        return env_value
    return default

def get_rclone_url() -> str:
    """
    Get the URL for the clone remote daemon
    """
    # Load default options
    host = get_parameter("RCLONE_HOST", "localhost")
    port = int(get_parameter("RCLONE_PORT", "5572"))
    protocol = get_parameter("RCLONE_PROTOCOL", "https")

    return f"{protocol}://{host}:{port}"

def backup(src: str, dest: str) -> requests.Response:
    """
    Run an Rclone backup from the given src to the given destination
    """
    base_url = get_rclone_url()
    url = f"{base_url}/sync/sync"

    return requests.post(url, json = {
        "_async": True,
        "_filter": {
            "FilterFrom": [ "/data/filter" ]
        },
        "srcFs": src,
        "dstFs": dest,
        "opt": {
            "retries": 1
        }
    })