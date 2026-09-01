import socket
import os
import time
import shutil
import json
from datetime import datetime



def record_check(name, status, details=""):
    """
    Save a check result to the logs list and print it nicely in the terminal.

    Parameters:
        name (str): Name of the check, e.g. "CPU Load"
        status (str): Result status, e.g. "OK" or "FAIL"
        details (str): Optional extra information
    """
    entry = {
        "timestamp": datetime.now().strftime("%Y-%m-%dT%H:%M:%S"),
        "check": name,
        "status": status,
        "details": details
    }
    logs.append(entry)

    # terminal colors
    color = "\033[92m" if status == "OK" else "\033[91m"
    reset = "\033[0m"
    print(f"{name:28} | {color}{status:8}{reset} | {details}")


def get_cpu_load():
    """
    Read system load averages from /proc/loadavg.

    Returns:
        str: load average for 1, 5, and 15 minutes
    """
    # On Linux, /proc/loadavg contains load averages:
    # e.g. "0.12 0.08 0.05 1/123 4567"
    with open("/proc/loadavg", "r") as f:
        load = f.read().split()[:3]
        return f"{load[0]} (1min), {load[1]} (5min), {load[2]} (15min)"


def get_local_ip():
    """
    Detect the machine's local IP address by creating a UDP socket.

    Returns:
        str: local IP address or 'Unknown'
    """
    s = None
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        # No actual data is sent; this just helps the OS choose a local interface
        s.connect(("8.8.8.8", 80))
        return s.getsockname()[0]
    except Exception:
        return "Unknown"
    finally:
        if s is not None:
            s.close()


def get_gateway():
    """
    Find the default gateway by parsing the output of `ip route`.

    Returns:
        str: gateway IP address or 'Unknown'
    """
    try:
        route_output = os.popen("ip route").read()
        for line in route_output.splitlines():
            if line.startswith("default"):
                return line.split()[2]
    except Exception:
        pass
    return "Unknown"


def test_tcp_connect(host, port, timeout=2):
    """
    Test whether a TCP connection can be established to host:port.

    Parameters:
        host (str): target hostname or IP
        port (int): target TCP port
        timeout (int|float): connection timeout in seconds

    Returns:
        tuple[bool, str]: (success, error_message)
    """
    s = None
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.settimeout(timeout)
        s.connect((host, port))
        return True, ""
    except Exception as e:
        return False, str(e)
    finally:
        if s is not None:
            s.close()



# Assuming the functions record_check, get_cpu_load, get_local_ip, 
# get_gateway, and test_tcp_connect are defined as in the previous response.

if __name__ == "__main__":
    # Print a clean, formatted header for the terminal output
    print("\n" + "="*60)
    print("Starting Network Health Diagnostics...")
    print("="*60)
    # The :28 and :8 are format specifiers to ensure fixed-width columns
    print(f"{'Check Name':28} | {'Status':8} | {'Details'}")
    print("-" * 75)

    # 1. System-level Checks
    # Get the machine's primary local IP address
    local_ip = get_local_ip()
    record_check("Local IP", "OK", local_ip)
    
    # Get the default gateway IP; flag as FAILED if not found
    gateway = get_gateway()
    record_check("Gateway", "OK" if gateway != "Unknown" else "FAILED", gateway)
    
    # Calculate free disk space on root directory (/) in Gigabytes
    # shutil.disk_usage returns (total, used, free); we take free and convert bytes to GB
    disk_free_gb = shutil.disk_usage('/').free // (2**30)
    record_check("Disk Space", "OK", f"{disk_free_gb} GB Free")

    # 2. Connectivity Checks
    # Define targets as a list of tuples: (name, host, port)
    tests = [
        ("Google DNS", "8.8.8.8", 53),
        ("Cloudflare DNS", "1.1.1.1", 53),
        ("Hmrah Academy", "hamrah.academy", 443),
    ]

    # Iterate through the list and perform the TCP connection test for each
    for name, host, port in tests:
        success, error = test_tcp_connect(host, port)
        status_str = "OK" if success else "FAILED"
        # Display the result to terminal and append to the logs list
        record_check(name, status_str, "Connected" if success else f"Error: {error}")

    # 3. Report Generation (Saving to JSON)
    # The 'logs' variable is populated via the record_check calls
    report_file = "health_report.json"
    try:
        with open(report_file, "w") as f:
            json.dump(logs, f, indent=4)
    except Exception as e:
        print(f"Failed to write report file: {e}")

    # Log CPU metrics last to capture the state after running other tests
    record_check("CPU Load", "OK", get_cpu_load())
    
    # Final footer to confirm completion
    print("-" * 75)
    print(f"Diagnostics Completed. Report saved to {report_file}")