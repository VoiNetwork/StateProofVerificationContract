from pyteal import *

"""
State Proof Verification Contract - clear state program
"""


def main():
    program = Seq([Approve()])

    return compileTeal(program, mode=Mode.Application, version=8)


print(main())
