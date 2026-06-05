#!/usr/bin/env python3
"""
Script to seed and delete error types from the Delícias da Lu API.

This script:
1. Authenticates with the API
2. Creates a list of error types
3. Optionally deletes them

Usage:
    python3 seed_error_types.py --create --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app --username admin --password password
    python3 seed_error_types.py --delete not-found --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app --username admin --password password
    python3 seed_error_types.py --create-and-delete --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app --username admin --password password
"""

import argparse
import json
import requests
import sys
from typing import Optional, List

# Color codes for terminal output
class Colors:
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    BLUE = '\033[94m'
    CYAN = '\033[96m'
    RESET = '\033[0m'
    BOLD = '\033[1m'

def print_info(msg: str):
    print(f"{Colors.BLUE}ℹ️  {msg}{Colors.RESET}")

def print_success(msg: str):
    print(f"{Colors.GREEN}✓ {msg}{Colors.RESET}")

def print_warning(msg: str):
    print(f"{Colors.YELLOW}⚠️  {msg}{Colors.RESET}")

def print_error(msg: str):
    print(f"{Colors.RED}✗ {msg}{Colors.RESET}")

def print_header(msg: str):
    print(f"\n{Colors.BOLD}{Colors.CYAN}{'='*60}{Colors.RESET}")
    print(f"{Colors.BOLD}{Colors.CYAN}{msg:^60}{Colors.RESET}")
    print(f"{Colors.BOLD}{Colors.CYAN}{'='*60}{Colors.RESET}\n")

# Default error types (from the API seed.go)
DEFAULT_ERROR_TYPES = [
    {
        "identifier": "not-found",
        "title": "Not Found",
        "detail": "The requested resource could not be found.",
        "status": 404,
    },
    {
        "identifier": "invalid-type",
        "title": "Invalid Type",
        "detail": "The provided type is not supported.",
        "status": 400,
    },
    {
        "identifier": "invalid-credentials",
        "title": "Invalid Credentials",
        "detail": "The username or password is incorrect.",
        "status": 401,
    },
    {
        "identifier": "invalid-token",
        "title": "Invalid Token",
        "detail": "The provided token is invalid or expired.",
        "status": 401,
    },
    {
        "identifier": "unauthorized",
        "title": "Unauthorized",
        "detail": "You are not authorized to access this resource.",
        "status": 401,
    },
    {
        "identifier": "type-not-found",
        "title": "Error Type Not Found",
        "detail": "No error type found for the provided identifier.",
        "status": 404,
    },
    {
        "identifier": "instance-not-found",
        "title": "Error Instance Not Found",
        "detail": "No error instance found for the provided identifier.",
        "status": 404,
    },
    {
        "identifier": "invalidDocumentLenght",
        "title": "Invalid Document Lenght",
        "detail": "The provided document does not have a known length (CPF or CNPJ).",
        "status": 400,
    },
    {
        "identifier": "invalidFilter",
        "title": "Invalid Filter",
        "detail": "The filter query parameter is invalid. Valid values are 'type' and 'instance'.",
        "status": 400,
    },
    {
        "identifier": "unexpectedUnhandledError",
        "title": "Unexpected Error",
        "detail": "An unexpected error occurred. Please contact support.",
        "status": 500,
    },
]

class APIClient:
    def __init__(self, base_url: str, username: str, password: str):
        self.base_url = base_url.rstrip('/')
        self.username = username
        self.password = password
        self.token: Optional[str] = None
        self.session = requests.Session()

    def authenticate(self) -> bool:
        """Authenticate with the API and get JWT token."""
        print_info(f"Authenticating as {self.username}...")
        
        url = f"{self.base_url}/v1/auth/login"
        payload = {
            "username": self.username,
            "password": self.password
        }
        
        try:
            response = self.session.post(url, json=payload, timeout=10)
            response.raise_for_status()
            
            data = response.json()
            self.token = data.get("token")
            
            if not self.token:
                print_error("No token in response")
                return False
            
            print_success(f"Authentication successful. Token: {self.token[:20]}...")
            return True
            
        except requests.exceptions.RequestException as e:
            print_error(f"Authentication failed: {e}")
            if hasattr(e, 'response') and e.response is not None:
                try:
                    print_error(f"Response: {e.response.json()}")
                except:
                    print_error(f"Response: {e.response.text}")
            return False

    def _get_headers(self) -> dict:
        """Get request headers with authentication."""
        headers = {
            "Content-Type": "application/json"
        }
        if self.token:
            headers["Authorization"] = f"Bearer {self.token}"
        return headers

    def create_error_types(self, error_types: List[dict]) -> bool:
        """Create error types via the API."""
        print_info(f"Creating {len(error_types)} error types...")
        
        url = f"{self.base_url}/v1/admin/error-types"
        
        try:
            response = self.session.post(
                url,
                json=error_types,
                headers=self._get_headers(),
                timeout=30
            )
            response.raise_for_status()
            
            data = response.json()
            print_success(f"Error types created: {data.get('message')}")
            print_success(f"Count: {data.get('count')}")
            return True
            
        except requests.exceptions.RequestException as e:
            print_error(f"Failed to create error types: {e}")
            if hasattr(e, 'response') and e.response is not None:
                try:
                    print_error(f"Response: {e.response.json()}")
                except:
                    print_error(f"Response: {e.response.text}")
            return False

    def delete_error_type(self, identifier: str) -> bool:
        """Delete a single error type."""
        url = f"{self.base_url}/v1/admin/error-types/{identifier}"
        
        try:
            response = self.session.delete(
                url,
                headers=self._get_headers(),
                timeout=10
            )
            response.raise_for_status()
            
            data = response.json()
            print_success(f"Deleted error type: {identifier}")
            return True
            
        except requests.exceptions.RequestException as e:
            print_error(f"Failed to delete error type '{identifier}': {e}")
            if hasattr(e, 'response') and e.response is not None:
                try:
                    print_error(f"Response: {e.response.json()}")
                except:
                    print_error(f"Response: {e.response.text}")
            return False

    def delete_all_error_types(self, error_types: List[dict]) -> bool:
        """Delete all error types."""
        print_info(f"Deleting {len(error_types)} error types...")
        
        failed = []
        for error_type in error_types:
            identifier = error_type.get("identifier")
            if not self.delete_error_type(identifier):
                failed.append(identifier)
        
        if failed:
            print_warning(f"Failed to delete {len(failed)} error types: {', '.join(failed)}")
            return False
        
        print_success(f"All error types deleted successfully")
        return True

def main():
    parser = argparse.ArgumentParser(
        description="Seed and delete error types for Delícias da Lu API"
    )
    
    parser.add_argument(
        "--api-url",
        default="https://delicias-da-lu-514609008596.southamerica-east1.run.app",
        help="API base URL"
    )
    parser.add_argument(
        "--username",
        required=True,
        help="API username"
    )
    parser.add_argument(
        "--password",
        required=True,
        help="API password"
    )
    
    # Action arguments (mutually exclusive)
    action_group = parser.add_mutually_exclusive_group(required=True)
    action_group.add_argument(
        "--create",
        action="store_true",
        help="Create error types"
    )
    action_group.add_argument(
        "--delete",
        type=str,
        help="Delete a specific error type by identifier"
    )
    action_group.add_argument(
        "--delete-all",
        action="store_true",
        help="Delete all error types"
    )
    action_group.add_argument(
        "--create-and-delete",
        action="store_true",
        help="Create error types, wait, then delete them"
    )
    
    parser.add_argument(
        "--wait",
        type=int,
        default=2,
        help="Wait time in seconds before deleting (only with --create-and-delete)"
    )
    
    args = parser.parse_args()
    
    # Initialize client
    client = APIClient(args.api_url, args.username, args.password)
    
    # Authenticate
    if not client.authenticate():
        sys.exit(1)
    
    print_header("Delícias da Lu - Error Types Management")
    
    # Perform action
    success = False
    
    if args.create:
        print_info(f"Creating {len(DEFAULT_ERROR_TYPES)} default error types...")
        success = client.create_error_types(DEFAULT_ERROR_TYPES)
        
    elif args.delete:
        print_info(f"Deleting error type: {args.delete}")
        success = client.delete_error_type(args.delete)
        
    elif args.delete_all:
        success = client.delete_all_error_types(DEFAULT_ERROR_TYPES)
        
    elif args.create_and_delete:
        print_info("Creating error types...")
        if not client.create_error_types(DEFAULT_ERROR_TYPES):
            sys.exit(1)
        
        print_info(f"Waiting {args.wait} seconds before deletion...")
        import time
        time.sleep(args.wait)
        
        success = client.delete_all_error_types(DEFAULT_ERROR_TYPES)
    
    # Print result
    print_header("Operation Complete")
    
    if success:
        print_success("All operations completed successfully")
        sys.exit(0)
    else:
        print_error("Some operations failed")
        sys.exit(1)

if __name__ == "__main__":
    main()
