# NSSM Usage Guide for Managing Rattler

## Introduction to NSSM

NSSM, short for Non-Sucking Service Manager, is a tool that allows you to easily create and manage Windows services for
any executable. It provides a simple interface to convert any executable into a Windows service, enabling you to start,
stop, and manage the service efficiently.

## Installation

1. **Download NSSM:**
    - Visit the NSSM GitHub releases page: [NSSM Releases](https://github.com/kirillkovalenko/nssm/releases).
    - Download the latest version of NSSM for your Windows architecture (32-bit or 64-bit).

2. **Extract NSSM:**
    - Extract the downloaded NSSM zip file to a location on your system.

3. **Add NSSM to System Path:**
    - Add the path to the NSSM executable (`nssm.exe`) to your system's PATH environment variable. This allows you to
      run NSSM from any directory in the command prompt.

## Deployment of Rattler with NSSM

1. **Register Rattler as a Service:**
    - Open a command prompt with administrator privileges.
    - Navigate to the directory where you extracted NSSM.
    - Execute the following command to register Rattler as a service:
      ```
      nssm install RattlerService "path\to\rattler.exe" -c ".rattler.yml"
      ```
   Replace `"path\to\rattler.exe"` with the actual path to your Rattler executable and `".rattler.yml"` with the
   configuration file parameter.

2. **Start the Rattler Service:**
    - Once Rattler is registered as a service, you can start it using the following command:
      ```
      nssm start RattlerService
      ```

3. **Stop the Rattler Service:**
    - To stop the Rattler service, use the following command:
      ```
      nssm stop RattlerService
      ```

4. **Delete the Rattler Service:**
    - If you no longer need the Rattler service, you can delete it with the following command:
      ```
      nssm remove RattlerService confirm
      ```

5. **Update the Rattler Service:**
    - If you need to update the configuration of the Rattler service, you can do so with the following steps:
        - Stop the Rattler service: `nssm stop RattlerService`
        - Modify the service configuration using NSSM GUI or command line options.
        - Start the Rattler service: `nssm start RattlerService`

## Managing Services with Task Manager

1. **Stop a Service:**
    - Press `Ctrl + Shift + Esc` to open Task Manager.
    - Switch to the "Services" tab.
    - Locate the Rattler service in the list.
    - Right-click on the service and select "Stop" from the context menu.

2. **Restart a Service:**
    - Follow the steps above to locate the Rattler service in Task Manager.
    - Right-click on the service and select "Restart" from the context menu.

## Conclusion

NSSM provides a convenient way to manage Windows services, including the deployment and management of Rattler as a
service. By following the steps outlined in this guide, you can easily integrate Rattler into your Windows environment
and ensure its smooth operation as a service. Additionally, Task Manager offers an alternative method for managing
services such as stopping and restarting them when necessary.