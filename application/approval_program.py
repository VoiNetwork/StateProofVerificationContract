from pyteal import *

"""
State Proof Verification Contract - approval program
"""


def main():
    # sets an admin address to global state
    on_creation = Seq(
        App.globalPut(Bytes("admin"), Txn.application_args[0]),
        Approve()
    )

    # checks if the interacting address is the admin address
    is_admin = Txn.sender() == App.globalGet(Bytes("admin"))

    program = Cond(
        [Txn.application_id() == Int(0), on_creation],
        # anyone can opt in
        [Txn.on_completion() == OnComplete.OptIn, Approve()],
        # only admin can delete or update
        [Txn.on_completion() == OnComplete.DeleteApplication, Return(is_admin)],
        [Txn.on_completion() == OnComplete.UpdateApplication, Return(is_admin)]
    )

    return compileTeal(
        program,
        mode=Mode.Application,
        version=8
    )


print(main())
