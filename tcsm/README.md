# Terraform Cloud State Manager

A very basic terraform cloud/enterprise state mager based on the following article:
<https://brendanthompson.com/posts/2021/rollback-terraform-state>

Currently it only supports state restoration and takes 3 mandatory flags as described below:

* Token: Your API token for Terraform Cloud
* Workspace ID: Your oraganisation's workspace id. That you can find on the workspace homepage in the `ws-<16 characters long alpha numeric ID>` format.
* State ID: The ID of the state you want to restore. You will find the ID mentioned with the state in the format `sv-<16 characters long alpha numeric ID>`.
