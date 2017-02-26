
Some of my learned knowledge written down, for myself.

## Notes
* The DPT is sometimes in the manufacturer data, which is in the "M-XXX" folders.
* The project is in the "P-XXXX" folder.
* Custom DPT exist, and the user can provide a configuration.

## Parameter of group objects
Basic Device Object configuration is stored in the `ComObjectTable>ComObject` inside the manufacturer data (this is the parameter stored in memory).
It's properties can be overwritten by `ComObjectRefs>ComObjectRef` (linked via RefId). 1 `ComObject` can have multiple `ComObjectRef`.
> A particular CO could need, depending on the “situation”, different DPTs.
ComObjectRefRef is used to visualize the parameters in ETS.

For more lookup "05 KNX Certification of Products - Procedure v01.03.01 AS.pdf 6.1.2.7.1.2"

Manual changes by the user overwrite all the above, if a user changes settings it is stored in 0.xml. 
For example custom Flags for device. Or custom DPT for address.

## Function name
It seems in newer ETS project the Function Name was removed.
TODO: verify this assumption

## Devices
Device are linked with the `0.xml` to the manufacturer XML. The DeviceInstance contains the device objects (ComObjectInstanceRefs).
Channels are grouped together and represent a hardware feature.

## Hardware identification
It seems `Hardware>Product.ID` (ProductRefId) is the same between all devices. The `Product.Hash` seems to be the same too.
And serial number of course should be the same.

## Bug: ETS shows no DPT but we have it
If ComObject defines a DPT but ComObjectRef overwrites it with empty "", then the DPT does not show in ETS.
I don't know why it would not show.

## Questions
* Hide unconnected Group Address? -> YES
* Can projects have multiple projects folders? (P-XXXX) -> no
* Can projects have multiple installations? (In 0.xml multiple "Installation") -> no
* Can a device objects channel be changed by the user in ETS?
* What are DeviceInstance>ParameterInstanceRef for? -> I guess they are device settings.
* Why does ComObjectRef sometimes have the same data as ComObject?
* If groupaddress has dpt and comobject has dpt, which one should I use?
* How do you set the ETS name property? -> ETS uses description and name only when not empty. Fallback is productname.
* Is the user allowed to change device Flags? -> It can be that it does not work when programming.
* What is the Send address on a device? -> Is the main address of a device. In a topology it is used to respond. The other addresses are just receive. Example: Taster only send values on the SEND address.
* What is the difference between ComObjectRef and ComObject?
* How is IsActive calculated? (red warning on icon) -> device requires plugin, maybe happend during upgrade.
* How can I find out if a configuration is written to the device. -> green checkmarks in topology view
* Changing a flag or removing an address removes the Grp tick. -> Grp is CommunicationPartLoaded
* Changing device parameter removes the Par tick. -> Par is ParametersLoaded
* Cfg is MediumConfigLoaded, no idea what it does.

