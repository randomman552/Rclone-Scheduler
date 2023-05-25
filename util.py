import requests
import os
from typing import Union

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

def apply_filter(body: dict) -> dict:
    """
    Apply global filter settings to the given request
    :param dict request: The request dict to edit
    """
    filter_from = get_parameter("BACKUP_FILTER_FROM", "")

    if len(filter_from) > 0:
        body["_filter"] = {
            "FilterFrom": [ filter_from ]
        }

    return body

def backup(src: str, dest: str) -> requests.Response:
    """
    Run an Rclone backup from the given src to the given destination
    """
    base_url = get_rclone_url()
    url = f"{base_url}/sync/sync"

    body = {
        "_async": True,
        "srcFs": src,
        "dstFs": dest,
        "opt": {
            "retries": 1
        }
    }

    body = apply_filter(body)

    return requests.post(url, json = body)