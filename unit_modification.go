package main

type unitModification struct {
	armor                int32
	magicResistance      int32
	criticalStrikeChance int32
	criticalStrikeDamage int32
	cooldownReduction    int32
	damageThreatFactor   int32
	healingThreatFactor  int32
}

// add adds two unitModifications
func (um *unitModification) add(operand *unitModification) {
	um.armor += operand.armor
	um.magicResistance += operand.magicResistance
	um.criticalStrikeChance += operand.criticalStrikeChance
	um.criticalStrikeDamage += operand.criticalStrikeDamage
	um.cooldownReduction += operand.cooldownReduction
	um.damageThreatFactor += operand.damageThreatFactor
	um.healingThreatFactor += operand.healingThreatFactor
}
