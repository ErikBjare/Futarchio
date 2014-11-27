# -*- coding: utf-8 -*-
from selenium import selenium
import unittest, time, re

class login(unittest.TestCase):
    def setUp(self):
        self.verificationErrors = []
        self.selenium = selenium("localhost", 4444, "*chrome", "http://localhost:8080/")
        self.selenium.start()
    
    def test_login(self):
        sel = self.selenium
        sel.open("/")
        sel.click("link=Login")
        time.sleep(0.5)
        sel.type("name=username", "erb")
        sel.type("name=password", "password")
        sel.click("//button[@type='submit']")
        sel.wait_for_page_to_load("30000")
        sel.click("link=My Account")
        sel.click("link=Log out")
        time.sleep(0.5)
        sel.click("link=Login")
        sel.click("css=span.title")
    
    def tearDown(self):
        self.selenium.stop()
        self.assertEqual([], self.verificationErrors)

if __name__ == "__main__":
    unittest.main()
