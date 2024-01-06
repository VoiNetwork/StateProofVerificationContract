from pyteal import *

"""
State Proof Verification Contract - approval program
"""


def main():
    # checks if the interacting address is the admin address
    is_admin = Txn.sender() == App.globalGet(Bytes("admin"))

    # sets an admin address to global state
    on_creation = Seq(
        App.globalPut(Bytes("admin"), Txn.application_args[0]),
        Approve()
    )

    # checks if the entry exists, if it does, it simply replaces it, otherwise it adds a new entry
    # on_add_state_proof = Seq(
    #     Assert(is_admin),
    #     entry := App.box_get(Txn.application_args[1]),
    #     If(entry.hasValue())
    #     .Then(App.box_replace(Txn.application_args[1], Int(0), Txn.application_args[2]))
    #     .Else(App.box_put(Txn.application_args[1], Txn.application_args[2])),
    #     Approve()
    # )

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
