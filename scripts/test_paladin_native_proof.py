"""Parent ownership cannot be rolled back or forged by appending a second section."""
import unittest
from paladin_native_proof import validate_parent


class ParentPrefixTests(unittest.TestCase):
    prefix='Accepted immutable pre-Paladin body.\n'
    section='\n<!-- edda-327:start -->\nPaladin #845.\n<!-- edda-327:end -->'

    def test_exact_owned_section_and_future_unrelated_append(self):
        validate_parent(self.prefix,self.section,self.prefix+self.section)
        validate_parent(self.prefix,self.section,self.prefix+self.section+'\nUnrelated later class section.')

    def test_rollback_duplicate_marker_and_altered_second_section(self):
        for actual in [self.prefix,self.prefix+self.section+self.section,self.prefix+self.section+'\n<!-- edda-327:start -->altered<!-- edda-327:end -->',self.prefix+self.section.replace('edda-327:end','edda-327:changed'),'Changed '+self.prefix+self.section]:
            with self.subTest(actual=actual),self.assertRaises(AssertionError):
                validate_parent(self.prefix,self.section,actual)


if __name__=='__main__':unittest.main()
