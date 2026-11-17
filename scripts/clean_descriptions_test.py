import unittest

from clean_descriptions import clean_doc

DOC = '''---
description: |-
  Use the *thing* resource.
  Example Usage

  data "x" "current" {}

  resource "thing" "a" {
    role_id = data.x.current.id

    user_ids = [data.x.current.id]
  }

  Permission reconciliation
  Only declared permissions are tracked.
---
'''


class CleanDocTest(unittest.TestCase):
    def test_drops_hcl_keeps_trailing_prose(self):
        self.assertEqual(
            clean_doc(DOC),
            '---\ndescription: |-\n  Use the &#42;thing&#42; resource. '
            'Permission reconciliation Only declared permissions are tracked.\n'
            '---\n')

    def test_rerun_is_noop(self):
        once = clean_doc(DOC)
        self.assertEqual(clean_doc(once), once)


if __name__ == '__main__':
    unittest.main()
