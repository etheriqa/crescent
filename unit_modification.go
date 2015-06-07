package main

type unitModification struct {
	armor                statistic
	magicResistance      statistic
	criticalStrikeChance statistic
	criticalStrikeFactor statistic
	cooldownReduction    statistic
	damageThreatFactor   statistic
	healingThreatFactor  statistic
}

// add adds two unitModifications
func (um *unitModification) add(operand *unitModification) {
	um.armor += operand.armor
	um.magicResistance += operand.magicResistance
	um.criticalStrikeChance += operand.criticalStrikeChance
	um.criticalStrikeFactor += operand.criticalStrikeFactor
	um.cooldownReduction += operand.cooldownReduction
	um.damageThreatFactor += operand.damageThreatFactor
	um.healingThreatFactor += operand.healingThreatFactor
}
